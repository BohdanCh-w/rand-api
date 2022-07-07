package integer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/bohdanch-w/rand-api/randapi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/urfave/cli/v2"
)

type coinParams struct {
	Format string
	Number int
}

func (p *coinParams) retriveParams(ctx *cli.Context) error {
	p.Format = ctx.String("format")
	p.Number = ctx.Int("number")

	return p.validate()
}

func (p *coinParams) validate() error {
	if err := validation.Validate(p.Format, validation.In("ukr", "eng", "number")); err != nil {
		return fmt.Errorf("`format` param is invalid: %w", err)
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

func Coin(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params coinParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		coinReq := integerRequest{
			ApiKey:      cfg.APIKey,
			Number:      params.Number,
			Min:         0,
			Max:         1,
			Replacement: true,
			Base:        10,
			PregenRand:  nil,
		}

		req, err := randapi.NewRandomRequest(method, coinReq)
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

		mapper, err := coinMappers(params.Format)
		if err != nil {
			return err
		}

		outputData := make([]interface{}, 0, len(data))
		for _, v := range data {
			side, err := coinSide(v, mapper)
			if err != nil {
				return fmt.Errorf("decode random data: %w", err)
			}

			outputData = append(outputData, side)
		}

		output.GenerateOutput(cfg.Output, outputData, apiInfo)

		return nil
	}
}

type coinMapper []interface{}

func coinMappers(format string) (coinMapper, error) {
	mapper, ok := map[string]coinMapper{
		"eng":    {"heads", "tails"},
		"ukr":    {"решка", "орел"},
		"number": {0, 1},
	}[format]

	if !ok {
		return nil, fmt.Errorf("invalid coin format: %s", format)
	}

	return mapper, nil
}

func coinSide(v int, mapper coinMapper) (interface{}, error) {
	if v < 0 || v > 1 {
		return nil, fmt.Errorf("invalid coin value %d", v)
	}

	return mapper[v], nil
}
