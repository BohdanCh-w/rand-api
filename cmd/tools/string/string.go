package string

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

const (
	method        = "generateStrings"
	maxStringLen  = 32
	maxCharsetLen = 128
	numberMax     = 10_000
)

type stringParams struct {
	Length  int
	Charset string
	Number  int
	Unique  bool
}

func (p *stringParams) retriveParams(ctx *cli.Context) error {
	p.Length = ctx.Int("length")
	p.Charset = ctx.String("charset")
	p.Number = ctx.Int("number")
	p.Unique = ctx.Bool("unique")

	return p.validate()
}

func (p *stringParams) validate() error {
	if err := validation.Validate(
		p.Length,
		validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(maxStringLen),
	); err != nil {
		return fmt.Errorf("`length` param is invalid: %w", err)
	}

	if err := validation.Validate(
		len(p.Charset),
		validation.Required.Error("length must be no less than 1"),
		validation.Max(maxCharsetLen),
	); err != nil {
		return fmt.Errorf("`charset` param is invalid: %w", err)
	}

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

func String(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params stringParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		strReq := stringRequest{
			ApiKey:      cfg.APIKey,
			Length:      params.Length,
			Characters:  params.Charset,
			Number:      params.Number,
			Replacement: !params.Unique,
			PregenRand:  nil,
		}

		req, err := randapi.NewRandomRequest(method, strReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := randapi.RandAPIExecute(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %w", err)
		}

		var (
			data    stringResponseData
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

type stringRequest struct {
	ApiKey      string  `json:"apiKey"`
	Length      int     `json:"length"`
	Characters  string  `json:"characters"`
	Number      int     `json:"n"`
	Replacement bool    `json:"replacement"`
	PregenRand  *string `json:"pregeneratedRandomization"`
}

type stringResponseData []string
