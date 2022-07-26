package gausian

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
)

const (
	CommandName     = "gausian"
	meanParam       = "mean"
	deviationParam  = "deviation"
	signDigitsParam = "signdig"
	numberParam     = "number"

	method               = "generateGaussians"
	rangeMaxMin          = 1_000_000
	minSignificantDigits = 2
	maxSignificantDigits = 14
	numberMax            = 10_000
)

// nolint: gomnd
func NewGausianCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:    CommandName,
		Aliases: []string{"gaus"},
		Usage:   "generate random value with Gausian distribution",
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:    meanParam,
				Usage:   "mean value of distribution         [-1000000, 1000000]",
				Aliases: []string{"m"},
				Value:   0,
			},
			&cli.Float64Flag{
				Name:    deviationParam,
				Usage:   "standart deviation of distribution [-1000000, 1000000]",
				Aliases: []string{"d"},
				Value:   1,
			},
			&cli.IntFlag{
				Name:    signDigitsParam,
				Usage:   "number of significant digits [2, 14]",
				Aliases: []string{"s"},
				Value:   6,
			},
			&cli.IntFlag{
				Name:    numberParam,
				Usage:   "number of values returned    [1, 10000]",
				Aliases: []string{"N"},
				Value:   1,
			},
		},
		Action: gausian(cfg),
	}
}

type gausianParams struct {
	Mean              float64
	Deviation         float64
	SignificantDigits int
	Number            int
}

func (p *gausianParams) retriveParams(ctx *cli.Context) error {
	p.Mean = ctx.Float64(meanParam)
	p.Deviation = ctx.Float64(deviationParam)
	p.SignificantDigits = ctx.Int(signDigitsParam)
	p.Number = ctx.Int(numberParam)

	return p.validate()
}

func (p *gausianParams) validate() error {
	if err := validation.Validate(
		p.Mean,
		validation.Min(float64(-rangeMaxMin)),
		validation.Max(float64(rangeMaxMin)),
	); err != nil {
		return fmt.Errorf("`mean` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Deviation,
		validation.Min(float64(-rangeMaxMin)),
		validation.Max(float64(rangeMaxMin)),
	); err != nil {
		return fmt.Errorf("`deviation` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.SignificantDigits,
		validation.Required.Error("must be no less than 2"),
		validation.Min(minSignificantDigits),
		validation.Max(maxSignificantDigits),
	); err != nil {
		return fmt.Errorf("`signdig` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Number, validation.Min(1),
		validation.Max(numberMax),
		validation.Required.Error("must be no less than 1"),
	); err != nil {
		return fmt.Errorf("`number` param is invalid: %w", err)
	}

	return nil
}

func gausian(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params gausianParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		gausReq := gausianRequest{
			APIKey:            cfg.APIKey,
			Mean:              params.Mean,
			StandardDeviation: params.Deviation,
			SignificantDigits: params.SignificantDigits,
			Number:            params.Number,
			PregenRand:        nil,
		}

		req, err := cfg.RandRetriever.NewRequest(method, gausReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %w", err)
		}

		var (
			data    gausianResponseData
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

		if err := cfg.OutputProcessor.GenerateRandOutput(outputData, apiInfo); err != nil {
			return fmt.Errorf("generate rand output: %w", err)
		}

		return nil
	}
}

type gausianRequest struct {
	APIKey            string  `json:"apiKey"`
	Mean              float64 `json:"mean"`
	StandardDeviation float64 `json:"standardDeviation"`
	SignificantDigits int     `json:"significantDigits"`
	Number            int     `json:"n"`
	PregenRand        *string `json:"pregeneratedRandomization"`
}

type gausianResponseData []float64
