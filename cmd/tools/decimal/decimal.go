package decimal

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
)

const (
	CommandName = "decimal"
	baseParam   = "base"
	placesParam = "places"
	numberParam = "number"
	uniqueParam = "unique"

	method           = "generateDecimalFractions"
	numberMax        = 10_000
	decimalPlacesMax = 14
)

// nolint: gomnd
func NewDecimalCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:    CommandName,
		Usage:   "generate random decimal value in range [0, 1]",
		Aliases: []string{"dec"},
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:    baseParam,
				Usage:   "returned value will be in range [0, base]",
				Aliases: []string{"b"},
				Value:   1,
			},
			&cli.IntFlag{
				Name:    placesParam,
				Usage:   "number of decimal places to use [1, 14]",
				Aliases: []string{"p"},
				Value:   6,
			},
			&cli.IntFlag{
				Name:    numberParam,
				Usage:   "number of values returned     [1, 10000]",
				Aliases: []string{"N"},
				Value:   1,
			},
			&cli.BoolFlag{
				Name:    uniqueParam,
				Usage:   "specifies whether values must be unique. Returns error if N < (to - from)",
				Aliases: []string{"u"},
			},
		},
		Action: randDecimal(cfg),
	}
}

type decimalParams struct {
	Base   float64
	Places int
	Number int
	Unique bool
}

func (p *decimalParams) retrieveParams(ctx *cli.Context) error {
	p.Base = ctx.Float64(baseParam)
	p.Places = ctx.Int(placesParam)
	p.Number = ctx.Int(numberParam)
	p.Unique = ctx.Bool(uniqueParam)

	return p.validate()
}

func (p *decimalParams) validate() error {
	const errMaxUniqueRandomExceeded = entities.Error("`number` of unique requested values is greater than possible")

	if err := validation.Validate(p.Base, validation.Required); err != nil {
		return fmt.Errorf("`base` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Places, validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(decimalPlacesMax),
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
		return fmt.Errorf("%w decimal places = %d", errMaxUniqueRandomExceeded, p.Places)
	}

	return nil
}

func randDecimal(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params decimalParams

		if err := params.retrieveParams(cCtx); err != nil {
			return err
		}

		decReq := decimalRequest{
			APIKey:        cfg.APIKey,
			Number:        params.Number,
			DecimalPlaces: params.Places,
			Replacement:   !params.Unique,
			PregenRand:    cfg.PregenRand,
		}

		req, err := cfg.RandRetriever.NewRequest(method, decReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
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

		if err := cfg.OutputProcessor.GenerateRandOutput(outputData, apiInfo); err != nil {
			return fmt.Errorf("generate rand output: %w", err)
		}

		return nil
	}
}

type decimalRequest struct {
	APIKey        string              `json:"apiKey"`
	Number        int                 `json:"n"`
	DecimalPlaces int                 `json:"decimalPlaces"`
	Replacement   bool                `json:"replacement"`
	PregenRand    entities.PregenRand `json:"pregeneratedRandomization"`
}

type decimalResponseData []float64
