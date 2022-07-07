package blob

import (
	"context"
	"encoding/base64"
	"encoding/hex"
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
	method    = "generateBlobs"
	sizeMax   = 1_048_576
	numberMax = 10_000

	hexFormat    = "hex"
	base64Format = "base64"
)

type blobParams struct {
	Size   int64
	Number int
	Hex    bool
}

func (p *blobParams) retriveParams(ctx *cli.Context) error {
	p.Size = ctx.Int64("size")
	p.Number = ctx.Int("number")
	p.Hex = ctx.Bool("hex")

	return p.validate()
}

func (p *blobParams) validate() error {
	if err := validation.Validate(
		p.Size,
		validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(sizeMax),
	); err != nil {
		return fmt.Errorf("`size` param is invalid: %w", err)
	}

	if err := validation.Validate(
		p.Number,
		validation.Required.Error("must be no less than 1"),
		validation.Min(1),
		validation.Max(numberMax),
	); err != nil {
		return fmt.Errorf("`number` param is invalid: %w", err)
	}

	if p.Size%8 != 0 {
		return fmt.Errorf("`size` parameter must be divisible by 8")
	}

	if totalSize := p.Size * int64(p.Number); totalSize > sizeMax {
		return fmt.Errorf("Total size %d must not exceed %d", totalSize, sizeMax)
	}

	return nil
}

func BLOB(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params blobParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		intReq := blobRequest{
			ApiKey:     cfg.APIKey,
			Size:       params.Size,
			Number:     params.Number,
			Format:     blobFormat(params.Hex),
			PregenRand: nil,
		}

		req, err := randapi.NewRandomRequest(method, intReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := randapi.RandAPIExecute(ctx, &req)
		if err != nil {
			return fmt.Errorf("get result: %w", err)
		}

		var (
			data    blobResponseData
			apiInfo = entities.APIInfo{
				ID:           req.ID,
				Timestamp:    time.Time(result.Random.Timestamp),
				RequestsLeft: result.RequestsLeft,
				BitsUsed:     result.BitsUsed,
				BitsLeft:     result.BitsLeft,
			}
			decoder = getDecoder(params.Hex)
		)

		if err := json.Unmarshal(result.Random.Data, &data); err != nil {
			return fmt.Errorf("decode result: %w", err)
		}

		outputData := make([]interface{}, 0, len(data))
		for _, v := range data {
			value, err := decoder(v)
			if err != nil {
				return fmt.Errorf("decode random data: %w", err)
			}

			outputData = append(outputData, string(value))
		}

		output.GenerateOutput(cfg.Output, outputData, apiInfo)

		return nil
	}
}

type blobRequest struct {
	ApiKey     string  `json:"apiKey"`
	Size       int64   `json:"size"`
	Number     int     `json:"n"`
	Format     string  `json:"format"`
	PregenRand *string `json:"pregeneratedRandomization"`
}

type blobResponseData []string

func blobFormat(isHex bool) string {
	if isHex {
		return hexFormat
	}

	return base64Format
}

type blobDecoder func(string) ([]byte, error)

func getDecoder(isHex bool) blobDecoder {
	if isHex {
		return hex.DecodeString
	}

	return base64.StdEncoding.DecodeString
}
