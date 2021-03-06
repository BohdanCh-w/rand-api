package uuid

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
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

const (
	method    = "generateUUIDs"
	numberMax = 10_000
)

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

func UUID(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params uuidParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		uuidReq := uuidRequest{
			ApiKey:     cfg.APIKey,
			Number:     params.Number,
			PregenRand: nil,
		}

		req, err := randapi.NewRandomRequest(method, uuidReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := randapi.RandAPIExecute(ctx, &req)
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

		output.GenerateOutput(cfg.Output, outputData, apiInfo)

		return nil
	}
}

type uuidRequest struct {
	ApiKey     string  `json:"apiKey"`
	Number     int     `json:"n"`
	PregenRand *string `json:"pregeneratedRandomization"`
}

type uuidResponseData []uuid.UUID
