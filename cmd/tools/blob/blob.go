package blob

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
)

const (
	CommandName = "blob"
	sizeParam   = "size"
	hexParam    = "hex"
	numberParam = "number"

	method    = "generateBlobs"
	sizeMax   = 1_048_576
	numberMax = 10_000

	hexFormat    = "hex"
	base64Format = "base64"
)

// nolint: gomnd
func NewBlobCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:  CommandName,
		Usage: "generate random Binary Large OBject. Total size must not exceed 1,048,576 bits (128 Kib)",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    sizeParam,
				Usage:   "size of blobs in bits [1, 1048576] must be divisible by 8",
				Aliases: []string{"s"},
				Value:   64,
			},
			&cli.BoolFlag{
				Name:        hexParam,
				Usage:       "if true generated data has hex format, base64 otherwise",
				DefaultText: base64Format,
			},
			&cli.IntFlag{
				Name:    numberParam,
				Usage:   "number of values returned [1, 10000]",
				Aliases: []string{"N"},
				Value:   1,
			},
		},
		Action: blob(cfg),
	}
}

type blobParams struct {
	Size   int64
	Number int
	Hex    bool
}

func (p *blobParams) retriveParams(ctx *cli.Context) error {
	p.Size = ctx.Int64(sizeParam)
	p.Number = ctx.Int(numberParam)
	p.Hex = ctx.Bool(hexParam)

	return p.validate()
}

func (p *blobParams) validate() error {
	const (
		errSizeDivisibleBy8 = entities.Error("`size` parameter must be divisible by 8")
		errSizeExceeded     = entities.Error("size exceed")
	)

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
		return errSizeDivisibleBy8
	}

	if totalSize := p.Size * int64(p.Number); totalSize > sizeMax {
		return fmt.Errorf("%w: %d > %d", errSizeExceeded, totalSize, sizeMax)
	}

	return nil
}

func blob(cfg *config.AppConfig) cli.ActionFunc { // nolint: funlen
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		var params blobParams

		if err := params.retriveParams(cCtx); err != nil {
			return err
		}

		blobReq := blobRequest{
			APIKey:     cfg.APIKey,
			Size:       params.Size,
			Number:     params.Number,
			Format:     blobFormat(params.Hex),
			PregenRand: cfg.PregenRand,
		}

		req, err := cfg.RandRetriever.NewRequest(method, blobReq)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}

		result, err := cfg.RandRetriever.ExecuteRequest(ctx, &req)
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

		if err := cfg.OutputProcessor.GenerateRandOutput(outputData, apiInfo); err != nil {
			return fmt.Errorf("generate rand output: %w", err)
		}

		return nil
	}
}

type blobRequest struct {
	APIKey     string              `json:"apiKey"`
	Size       int64               `json:"size"`
	Number     int                 `json:"n"`
	Format     string              `json:"format"`
	PregenRand entities.PregenRand `json:"pregeneratedRandomization"`
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
