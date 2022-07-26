package uuid

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

const (
	method    = "generateUUIDs"
	numberMax = 10_000
)

func NewUUIDCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:  "uuid",
		Usage: "generate random uuid V4",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "number",
				Usage:   "number of values returned [1, 10000]",
				Aliases: []string{"N"},
				Value:   1,
			},
		},
		Action: randUUID(cfg),
	}
}

type uuidParams struct {
	Number int
}

func (p *uuidParams) retriveParams(ctx *cli.Context) error {
	p.Number = ctx.Int("number")

	return p.validate()
}

func (p *uuidParams) validate() error {
	if err := validation.Validate(
		p.Number,
		validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(numberMax),
	); err != nil {
		return fmt.Errorf("`number` param is invalid: %w", err)
	}

	return nil
}

func randUUID(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params uuidParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		uuidReq := uuidRequest{
			APIKey:     cfg.APIKey,
			Number:     params.Number,
			PregenRand: nil,
		}

		req, err := cfg.RandRetriever.NewRequest(method, uuidReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %w", err)
		}

		var (
			data    uuidResponseData
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

type uuidRequest struct {
	APIKey     string  `json:"apiKey"`
	Number     int     `json:"n"`
	PregenRand *string `json:"pregeneratedRandomization"`
}

type uuidResponseData []uuid.UUID
