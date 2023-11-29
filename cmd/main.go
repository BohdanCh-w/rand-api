package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bohdanch-w/rand-api/cmd/tools/blob"
	"github.com/bohdanch-w/rand-api/cmd/tools/coin"
	"github.com/bohdanch-w/rand-api/cmd/tools/decimal"
	"github.com/bohdanch-w/rand-api/cmd/tools/gausian"
	"github.com/bohdanch-w/rand-api/cmd/tools/integer"
	"github.com/bohdanch-w/rand-api/cmd/tools/status"
	randstr "github.com/bohdanch-w/rand-api/cmd/tools/string"
	"github.com/bohdanch-w/rand-api/cmd/tools/uuid"
	"github.com/bohdanch-w/rand-api/cmd/tools/version"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/bohdanch-w/rand-api/randapi"

	guuid "github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

//go:embed resources/*
var apiKeyResource embed.FS

const (
	CommandName     = "randapi"
	apiPathParam    = "api-path"
	apikeyParam     = "apikey"
	pregenIDParam   = "pr-id"
	pregenDateParam = "pr-date"
	signedParam     = "signed"
	verboseParam    = "verbose"
	quietParam      = "quiet"
	timeoutParam    = "timeout"
	separatorParam  = "separator"
	outputParam     = "file"

	defaultTimeout     = 5 * time.Second
	defaultSeparator   = " "
	defaultRandAPIPath = "https://api.random.org/json-rpc/4/invoke"
)

func retriveParamsFunc(cfg *config.AppConfig, f **os.File) cli.BeforeFunc {
	return func(c *cli.Context) error {
		const (
			errInvalidDate         = entities.Error("latest allowed date is today")
			errBothPregenSpecified = entities.Error("only pr-id OR pr-date is allowed. Not both")
		)

		if apiKey := c.String(apikeyParam); len(apiKey) != 0 {
			if _, err := guuid.Parse(apiKey); err != nil {
				return fmt.Errorf("api-key: %w", err)
			}

			cfg.APIKey = apiKey
		}

		if prDate := c.String(pregenDateParam); len(prDate) > 0 {
			d, err := time.Parse("2006-01-02", prDate)
			if err != nil {
				return fmt.Errorf("invlid date format: %w", err)
			}

			if d.After(time.Now().UTC()) {
				return errInvalidDate
			}

			cfg.PregenRand.Date = &prDate
		}

		if prID := c.String(pregenIDParam); len(prID) > 0 {
			cfg.PregenRand.ID = &prID
		}

		if cfg.PregenRand.ID != nil && cfg.PregenRand.Date != nil {
			return errBothPregenSpecified
		}

		cfg.Timeout = c.Duration(timeoutParam)

		var (
			w   io.Writer = os.Stdout
			err error
		)

		if destination := c.String(outputParam); destination != "" {
			*f, err = os.Create(destination)
			if err != nil {
				return fmt.Errorf("create output file: %w", err)
			}

			w = *f
		}

		cfg.OutputProcessor = output.NewOutputProcessor(
			c.Bool(verboseParam),
			c.Bool(quietParam),
			c.String(separatorParam),
			w,
		)

		cfg.RandRetriever = randapi.NewRandomOrgRetriever(
			c.String(apiPathParam),
			http.DefaultClient,
			c.Bool(signedParam),
		)

		return nil
	}
}

func main() { // nolint: funlen
	var (
		apiKeyRequired = true
		cfg            config.AppConfig
		apiKeyData, _  = apiKeyResource.ReadFile("resources/api-key")
		versionData, _ = apiKeyResource.ReadFile("resources/version")
		f              *os.File
	)

	if len(apiKeyData) > 0 {
		cfg.APIKey = string(apiKeyData)
		cfg.Version = string(versionData)
		apiKeyRequired = false
	}

	app := &cli.App{
		Name:  CommandName,
		Usage: "cli program to retrieve values from random.org",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    signedParam,
				Aliases: []string{"s"},
				Usage:   "get signed reply from random.org",
			},
			&cli.StringFlag{
				Name:  apiPathParam,
				Usage: "random api path",
				Value: defaultRandAPIPath,
			},
			&cli.StringFlag{
				Name:        apikeyParam,
				Usage:       "specify custom apikey",
				DefaultText: "embedded resource",
				Required:    apiKeyRequired,
			},
			&cli.StringFlag{
				Name:  pregenIDParam,
				Usage: "pregenerated randomization by id string",
			},
			&cli.StringFlag{
				Name:  pregenDateParam,
				Usage: "pregenerated randomization by date",
			},
			&cli.BoolFlag{
				Name:    verboseParam,
				Aliases: []string{"v"},
				Usage:   "make verbose output after completion",
			},
			&cli.BoolFlag{
				Name:    quietParam,
				Aliases: []string{"q"},
				Usage:   "suppress all warnings",
			},
			&cli.DurationFlag{
				Name:    timeoutParam,
				Aliases: []string{"t"},
				Usage:   "randomness server response timeout in seconds",
				Value:   defaultTimeout,
			},
			&cli.StringFlag{
				Name:    separatorParam,
				Aliases: []string{"sep"},
				Usage:   "string to separate output",
				Value:   defaultSeparator,
			},
			&cli.StringFlag{
				Name:    outputParam,
				Aliases: []string{"o"},
				Usage:   "save output to specified file",
			},
		},
		Before: retriveParamsFunc(&cfg, &f),
		Commands: []*cli.Command{
			integer.NewIntegerCommand(&cfg),
			coin.NewCoinCommand(&cfg),
			decimal.NewDecimalCommand(&cfg),
			gausian.NewGausianCommand(&cfg),
			randstr.NewStringCommand(&cfg),
			uuid.NewUUIDCommand(&cfg),
			blob.NewBlobCommand(&cfg),
			status.NewStatusCommand(&cfg),
			version.NewVersionCommand(&cfg),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
