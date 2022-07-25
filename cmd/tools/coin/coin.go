package coin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/urfave/cli/v2"
)

const (
	coinCommandName = "coin"
	formatParam     = "format"
	numberParam     = "number"

	method      = "generateIntegers"
	rangeMaxMin = 1_000_000_000
	numberMax   = 10_000

	formatEng = "eng"
	formatUkr = "ukr"
	formatNum = "num"
)

func NewCoinCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:  coinCommandName,
		Usage: "generate random coinflip result (two values possible)",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    numberParam,
				Usage:   "number of values returned [-1000, 1000]",
				Aliases: []string{"N"},
				Value:   1,
			},
			&cli.StringFlag{
				Name:    formatParam,
				Usage:   "format printet result. One of 'eng' 'ukr' 'num'",
				Aliases: []string{"f"},
				Value:   "eng",
			},
		},
		Action: coin(cfg),
	}
}

type coinParams struct {
	Format string
	Number int
}

func (p *coinParams) retriveParams(ctx *cli.Context) error {
	p.Format = ctx.String(formatParam)
	p.Number = ctx.Int(numberParam)

	return p.validate()
}

func (p *coinParams) validate() error {
	if err := validation.Validate(p.Format, validation.In(formatUkr, formatEng, formatNum)); err != nil {
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

func coin(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params coinParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		coinReq := coinRequest{
			ApiKey:      cfg.APIKey,
			Number:      params.Number,
			Min:         0,
			Max:         1,
			Replacement: true,
			Base:        10,
			PregenRand:  nil,
		}

		req, err := cfg.RandRetriever.NewRequest(method, coinReq)
		if err != nil {
			return fmt.Errorf("create request: %v", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %v", err)
		}

		var (
			data    coinResponseData
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

		cfg.OutputProcessor.GenerateRandOutput(outputData, apiInfo)

		return nil
	}
}

type coinMapper []interface{}

func coinMappers(format string) (coinMapper, error) {
	mapper, ok := map[string]coinMapper{
		formatEng: {"heads", "tails"},
		formatUkr: {"решка", "орел"},
		formatNum: {0, 1},
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

type coinRequest struct {
	ApiKey      string  `json:"apiKey"`
	Number      int     `json:"n"`
	Min         int64   `json:"min"`
	Max         int64   `json:"max"`
	Replacement bool    `json:"replacement"`
	Base        int8    `json:"base"`
	PregenRand  *string `json:"pregeneratedRandomization"`
}

type coinResponseData []int
