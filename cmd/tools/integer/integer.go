package integer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/bohdanch-w/rand-api/randapi"
)

const (
	method      = "generateIntegers"
	rangeMaxMin = 1_000_000_000
	numberMax   = 10_000
)

type integerParams struct {
	From   int64
	To     int64
	Number int
	Unique bool
}

func retriveParams(ctx *cli.Context) (integerParams, error) {
	p := integerParams{
		From:   ctx.Int64("from"),
		To:     ctx.Int64("to"),
		Number: ctx.Int("number"),
		Unique: ctx.Bool("unique"),
	}

	if err := p.validate(); err != nil {
		return integerParams{}, fmt.Errorf("interger: %w", err)
	}

	return p, nil
}

func (p *integerParams) validate() error {
	if err := validation.Validate(p.From, validation.Min(-rangeMaxMin), validation.Max(rangeMaxMin)); err != nil {
		return fmt.Errorf("`from` param is invalid: %w", err)
	}

	if err := validation.Validate(p.To, validation.Min(-rangeMaxMin), validation.Max(rangeMaxMin)); err != nil {
		return fmt.Errorf("`to` param is invalid: %w", err)
	}

	if err := validation.Validate(p.Number, validation.Min(1), validation.Max(rangeMaxMin)); err != nil {
		return fmt.Errorf("`number` param is invalid: %w", err)
	}

	return nil
}

func Integer(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		params, err := retriveParams(cCtx)
		if err != nil {
			return err
		}

		intReq := integerRequest{
			ApiKey:      cfg.APIKey,
			Number:      params.Number,
			Min:         params.From,
			Max:         params.To,
			Replacement: params.Unique,
			Base:        10,
			PregenRand:  nil,
		}

		req, err := randapi.NewRandomRequest(method, intReq)
		if err != nil {
			return fmt.Errorf("create request: %v", err)
		}

		result, err := randapi.RandAPIExecute(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %v", err)
		}

		var (
			data    integerResponseData
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

		outputData := make([]interface{}, 0, len(data))
		for _, v := range data {
			outputData = append(outputData, v)
		}

		output.GenerateOutput(cfg.Output, outputData, apiInfo)

		return nil
	}
}

type integerRequest struct {
	ApiKey      string  `json:"apiKey"`
	Number      int     `json:"n"`
	Min         int64   `json:"min"`
	Max         int64   `json:"max"`
	Replacement bool    `json:"replacement"`
	Base        int8    `json:"base"`
	PregenRand  *string `json:"pregeneratedRandomization"`
}

type integerResponseData []int