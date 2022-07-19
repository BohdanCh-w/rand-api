package main

import (
	"embed"
	"fmt"
	"log"
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
	"github.com/bohdanch-w/rand-api/config"

	guuid "github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

//go:embed resources/*
var apiKeyResource embed.FS

const (
	CommandName    = "randapi"
	apikeyParam    = "apikey"
	signedParam    = "signed"
	verboseParam   = "verbose"
	quietParam     = "quiet"
	timeoutParam   = "timeout"
	separatorParam = "separator"
	outputParam    = "file"

	defaultTimeout   = 5 * time.Second
	defaultSeparator = " "
)

func retriveParamsFunc(cfg *config.AppConfig) cli.BeforeFunc {
	return func(c *cli.Context) error {
		if apiKey := c.String(apikeyParam); len(apiKey) != 0 {
			if _, err := guuid.Parse(apiKey); err != nil {
				return fmt.Errorf("api-key: %w", err)
			}

			cfg.APIKey = apiKey
		}

		if output := c.String(outputParam); output != "" {
			cfg.Output.Filename = &output
		}

		cfg.Signed = c.Bool(signedParam)
		cfg.Timeout = c.Duration(timeoutParam)
		cfg.Output.Verbose = c.Bool(verboseParam)
		cfg.Output.Quite = c.Bool(quietParam)
		cfg.Output.Separator = c.String(separatorParam)

		return nil
	}
}

func main() {
	var (
		apiKeyRequired = true
		cfg            config.AppConfig
		data, _        = apiKeyResource.ReadFile("resources/api-key")
	)

	if len(data) > 0 {
		cfg.APIKey = string(data)
		apiKeyRequired = false
	}

	app := &cli.App{
		Name:  "randapi",
		Usage: "cli program to retrieve values from random.org",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    signedParam,
				Aliases: []string{"s"},
				Usage:   "get signed reply from random.org",
			},
			&cli.StringFlag{
				Name:        apikeyParam,
				Usage:       "specify custom apikey",
				DefaultText: "embeded resource",
				Required:    apiKeyRequired,
			},
			&cli.BoolFlag{
				Name:    verboseParam,
				Aliases: []string{"v"},
				Usage:   "make verbose output after completition",
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
				Usage:   "save output to specied file",
			},
		},
		Before: retriveParamsFunc(&cfg),
		Commands: []*cli.Command{
			integer.NewIntegerCommand(&cfg),
			coin.NewCoinCommand(&cfg),
			decimal.NewDecimalCommand(&cfg),
			gausian.NewGausianCommand(&cfg),
			randstr.NewStringCommand(&cfg),
			uuid.NewUUIDCommand(&cfg),
			blob.NewBlobCommand(&cfg),
			status.NewStatusCommand(&cfg),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
