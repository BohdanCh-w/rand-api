package decimal

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/bohdanch-w/rand-api/randapi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli/v2"
)

const (
	method           = "generateDecimalFractions"
	numberMax        = 10_000
	decimalPlacesMax = 14
)

type decimalParams struct {
	Base   float64
	Places int
	Number int
	Unique bool
}

func (p *decimalParams) retrieveParams(ctx *cli.Context) error {
	p.Base = ctx.Float64("base")
	p.Places = ctx.Int("places")
	p.Number = ctx.Int("number")
	p.Unique = ctx.Bool("unique")

	return p.validate()
}

func (p *decimalParams) validate() error {
	if err := validation.Validate(p.Base, validation.Required); err != nil {
		return fmt.Errorf("`base` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Places, validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(14),
	); err != nil {
		return fmt.Errorf("`places` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Number, validation.Min(1),
		validation.Max(numberMax),
		validation.Required.Error("must be no less than 1"),
	); err != nil {
		return fmt.Errorf("`number` param is invalid: %w", err)
	}

	if p.Number > int(math.Pow10(p.Places)) {
		return fmt.Errorf("`number` of unique requested values is greater than possible with decimal places = %d", p.Places)
	}

	return nil
}

func Decimal(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params decimalParams

		if err := params.retrieveParams(cCtx); err != nil {
			return err
		}

		decReq := decimalRequest{
			ApiKey:        cfg.APIKey,
			Number:        params.Number,
			DecimalPlaces: params.Places,
			Replacement:   !params.Unique,
			PregenRand:    nil,
		}

		req, err := randapi.NewRandomRequest(method, decReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := randapi.RandAPIExecute(ctx, &req)
		if err != nil {
			return fmt.Errorf("get response: %w", err)
		}

		var (
			data    decimalResponseData
			apiInfo = entities.APIInfo{
				ID:           req.ID,
				Timestamp:    time.Time(result.Random.Timestamp),
				RequestsLeft: result.RequestsLeft,
				BitsUsed:     result.BitsUsed,
				BitsLeft:     result.BitsLeft,
			}
		)

		if err := json.Unmarshal(result.Random.Data, &data); err != nil {
			return fmt.Errorf("decode result: %w", err)
		}

		base := decimal.NewFromFloat(params.Base)

		outputData := make([]interface{}, 0, len(data))
		for _, v := range data {
			value, _ := base.Mul(decimal.NewFromFloat(v)).Float64()
			outputData = append(outputData, value)
		}

		output.GenerateOutput(cfg.Output, outputData, apiInfo)

		return nil
	}
}

type decimalRequest struct {
	ApiKey        string  `json:"apiKey"`
	Number        int     `json:"n"`
	DecimalPlaces int     `json:"decimalPlaces"`
	Replacement   bool    `json:"replacement"`
	PregenRand    *string `json:"pregeneratedRandomization"`
}

type decimalResponseData []float64
