package integer

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
	commandName = "integer"
	fromParam   = "from"
	toParam     = "to"
	numberParam = "number"
	uniqueParam = "unique"

	method      = "generateIntegers"
	rangeMaxMin = 1_000_000_000
	numberMax   = 10_000

	intBase = 10
)

// nolint: gomnd
func NewIntegerCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:    commandName,
		Aliases: []string{"int"},
		Usage:   "generate random integer in range (including)",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    fromParam,
				Usage:   "bottom limit of random number [-1e9, 1e9]",
				Aliases: []string{"f"},
				Value:   1,
			},
			&cli.Int64Flag{
				Name:    toParam,
				Usage:   "upper limit of random number  [-1e9, 1e9]",
				Aliases: []string{"t"},
				Value:   100,
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
		Action: integer(cfg),
	}
}

type integerParams struct {
	From   int64
	To     int64
	Number int
	Unique bool
}

func (p *integerParams) retriveParams(ctx *cli.Context) error {
	p.From = ctx.Int64(fromParam)
	p.To = ctx.Int64(toParam)
	p.Number = ctx.Int(numberParam)
	p.Unique = ctx.Bool(uniqueParam)

	return p.validate()
}

func (p *integerParams) validate() error {
	const (
		errToBiggerThanFrom        = entities.Error("`from` param must be less than `to`")
		errMaxUniqueRandomExceeded = entities.Error("`number` of unique requested values is greater than possible")
	)

	if err := validation.Validate(p.From, validation.Min(-rangeMaxMin), validation.Max(rangeMaxMin)); err != nil {
		return fmt.Errorf("`from` param is invalid: %w", err)
	}

	if err := validation.Validate(p.To, validation.Min(-rangeMaxMin), validation.Max(rangeMaxMin)); err != nil {
		return fmt.Errorf("`to` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Number,
		validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(numberMax),
	); err != nil {
		return fmt.Errorf("`number` param is invalid: %w", err)
	}

	if p.From >= p.To {
		return fmt.Errorf("%w: from %d to %d", errToBiggerThanFrom, p.From, p.To)
	}

	if (p.To-p.From) < int64(p.Number) && p.Unique {
		return fmt.Errorf("%w in range %d - %d", errMaxUniqueRandomExceeded, p.From, p.To)
	}

	return nil
}

func integer(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params integerParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		intReq := integerRequest{
			APIKey:      cfg.APIKey,
			Number:      params.Number,
			Min:         params.From,
			Max:         params.To,
			Replacement: !params.Unique,
			Base:        intBase,
			PregenRand:  cfg.PregenRand,
		}

		req, err := cfg.RandRetriever.NewRequest(method, intReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %w", err)
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

		if err := cfg.OutputProcessor.GenerateRandOutput(outputData, apiInfo); err != nil {
			return fmt.Errorf("generate rand output: %w", err)
		}

		return nil
	}
}

type integerRequest struct {
	APIKey      string              `json:"apiKey"`
	Number      int                 `json:"n"`
	Min         int64               `json:"min"`
	Max         int64               `json:"max"`
	Replacement bool                `json:"replacement"`
	Base        int8                `json:"base"`
	PregenRand  entities.PregenRand `json:"pregeneratedRandomization"`
}

type integerResponseData []int
