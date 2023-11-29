package coin

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
	coinCommandName = "coin"
	formatParam     = "format"
	numberParam     = "number"

	method    = "generateIntegers"
	numberMax = 10_000

	formatEng = "eng"
	formatUkr = "ukr"
	formatNum = "num"

	intBase = 10
)

const (
	errMapperInvalidFormat = entities.Error("coin mapper invalid format")
	errMapperInvalidValue  = entities.Error("coin mapper invalid value")
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
			APIKey:      cfg.APIKey,
			Number:      params.Number,
			Min:         0,
			Max:         1,
			Replacement: true,
			Base:        intBase,
			PregenRand:  cfg.PregenRand,
		}

		req, err := cfg.RandRetriever.NewRequest(method, coinReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %w", err)
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

		outputData, err := newFunction(params, data)
		if err != nil {
			return err
		}

		if err := cfg.OutputProcessor.GenerateRandOutput(outputData, apiInfo); err != nil {
			return fmt.Errorf("generate rand output: %w", err)
		}

		return nil
	}
}

func newFunction(params coinParams, data coinResponseData) ([]interface{}, error) {
	mapper, err := coinMappers(params.Format)
	if err != nil {
		return nil, err
	}

	outputData := make([]interface{}, 0, len(data))

	for _, v := range data {
		side, err := coinSide(v, mapper)
		if err != nil {
			return nil, fmt.Errorf("decode random data: %w", err)
		}

		outputData = append(outputData, side)
	}

	return outputData, nil
}

type coinMapper []interface{}

func coinMappers(format string) (coinMapper, error) {
	mapper, ok := map[string]coinMapper{
		formatEng: {"heads", "tails"},
		formatUkr: {"решка", "орел"},
		formatNum: {0, 1},
	}[format]

	if !ok {
		return nil, fmt.Errorf("%w: %s", errMapperInvalidFormat, format)
	}

	return mapper, nil
}

func coinSide(v int, mapper coinMapper) (interface{}, error) {
	if v < 0 || v > 1 {
		return nil, fmt.Errorf("%w: %d", errMapperInvalidValue, v)
	}

	return mapper[v], nil
}

type coinRequest struct {
	APIKey      string              `json:"apiKey"`
	Number      int                 `json:"n"`
	Min         int64               `json:"min"`
	Max         int64               `json:"max"`
	Replacement bool                `json:"replacement"`
	Base        int8                `json:"base"`
	PregenRand  entities.PregenRand `json:"pregeneratedRandomization"`
}

type coinResponseData []int
