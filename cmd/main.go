package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bohdanch-w/rand-api/cmd/tools/blob"
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

const characterRange = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

//go:embed resources/*
var apiKeyResource embed.FS

func main() {
	var (
		timeout        int
		signed         bool
		verbose        bool
		quite          bool
		apiKey         string
		output         string
		separator      string
		apiKeyRequired = true
		cfg            config.AppConfig
	)

	data, _ := apiKeyResource.ReadFile("resources/api-key")

	if len(data) > 0 {
		cfg.APIKey = string(data)
		apiKeyRequired = false
	}

	app := &cli.App{
		Name:  "rand",
		Usage: "cli program to retrieve values from random.org",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "signed",
				Aliases:     []string{"s"},
				Usage:       "get signed reply from random.org",
				Destination: &signed,
			},
			&cli.StringFlag{
				Name:        "apikey",
				Usage:       "specify custom apikey",
				DefaultText: "embeded resource",
				Required:    apiKeyRequired,
				Destination: &apiKey,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "make verbose output after completition",
				Destination: &verbose,
			},
			&cli.BoolFlag{
				Name:        "quite",
				Aliases:     []string{"q"},
				Usage:       "suppress all warnings",
				Destination: &quite,
			},
			&cli.IntFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Usage:       "randomness server response timeout in seconds",
				Value:       5,
				DefaultText: "5 seconds",
				Destination: &timeout,
			},
			&cli.StringFlag{
				Name:        "separator",
				Aliases:     []string{"sep"},
				Usage:       "string to separate output",
				Value:       " ",
				Destination: &separator,
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "save output to specied file",
				Destination: &output,
			},
		},
		Before: func(cCtx *cli.Context) error {
			if len(apiKey) != 0 {
				if _, err := guuid.Parse(apiKey); err != nil {
					return fmt.Errorf("api key: %w", err)
				}

				cfg.APIKey = apiKey
			}

			if output != "" {
				cfg.Output.Filename = &output
			}

			cfg.Signed = signed
			cfg.Timeout = time.Duration(timeout) * time.Second

			cfg.Output.Verbose = verbose
			cfg.Output.Quite = quite
			cfg.Output.Separator = separator

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "integer",
				Aliases: []string{"int"},
				Usage:   "generate random integer in range (including)",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "from",
						Usage:   "bottom limit of random number [-1e9, 1e9]",
						Aliases: []string{"f"},
						Value:   1,
					},
					&cli.Int64Flag{
						Name:    "to",
						Usage:   "upper limit of random number  [-1e9, 1e9]",
						Aliases: []string{"t"},
						Value:   100,
					},
					&cli.IntFlag{
						Name:    "number",
						Usage:   "number of values returned     [1, 10000]",
						Aliases: []string{"N"},
						Value:   1,
					},
					&cli.BoolFlag{
						Name:    "unique",
						Usage:   "specifies whether values must be unique. Returns error if N < (to - from)",
						Aliases: []string{"u"},
					},
				},
				Action: integer.Integer(&cfg),
			},
			{
				Name:  "coin",
				Usage: "generate random coinflip result (two values possible)",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "number",
						Usage:   "number of values returned [-1000, 1000]",
						Aliases: []string{"N"},
						Value:   1,
					},
					&cli.StringFlag{
						Name:    "format",
						Usage:   "format printet result. One of 'eng' 'ukr' 'number'",
						Aliases: []string{"f"},
						Value:   "eng",
					},
				},
				Action: integer.Coin(&cfg),
			},
			{
				Name:    "decimal",
				Usage:   "generate random decimal value in range [0, 1]",
				Aliases: []string{"dec"},
				Flags: []cli.Flag{
					&cli.Float64Flag{
						Name:    "base",
						Usage:   "returned value will be in range [0, base]",
						Aliases: []string{"b"},
						Value:   1,
					},
					&cli.IntFlag{
						Name:    "places",
						Usage:   "number of decimal places to use [1, 14]",
						Aliases: []string{"p"},
						Value:   6,
					},
					&cli.IntFlag{
						Name:    "number",
						Usage:   "number of values returned     [1, 10000]",
						Aliases: []string{"N"},
						Value:   1,
					},
					&cli.BoolFlag{
						Name:    "unique",
						Usage:   "specifies whether values must be unique. Returns error if N < (to - from)",
						Aliases: []string{"u"},
					},
				},
				Action: decimal.Decimal(&cfg),
			},
			{
				Name:    "gausian",
				Aliases: []string{"gaus"},
				Usage:   "generate random value with Gausian distribution",
				Flags: []cli.Flag{
					&cli.Float64Flag{
						Name:    "mean",
						Usage:   "mean value of distribution         [-1000000, 1000000]",
						Aliases: []string{"m"},
						Value:   0,
					},
					&cli.Float64Flag{
						Name:    "deviation",
						Usage:   "standart deviation of distribution [-1000000, 1000000]",
						Aliases: []string{"d"},
						Value:   1,
					},
					&cli.IntFlag{
						Name:    "signdig",
						Usage:   "number of significant digits [2, 14]",
						Aliases: []string{"s"},
						Value:   6,
					},
					&cli.IntFlag{
						Name:    "number",
						Usage:   "number of values returned    [1, 10000]",
						Aliases: []string{"N"},
						Value:   1,
					},
				},
				Action: gausian.Gausian(&cfg),
			},
			{
				Name:    "string",
				Usage:   "generate random string of given characters",
				Aliases: []string{"str"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "length",
						Usage:   "length of generated strings [1, 32](If N > 1, all strings have same length)",
						Aliases: []string{"l"},
						Value:   1,
					},
					&cli.StringFlag{
						Name:        "charset",
						Usage:       "characters to be used in generation. Max len - 128",
						Aliases:     []string{"c"},
						Value:       characterRange,
						DefaultText: "[A-Za-z0-9_]",
					},
					&cli.IntFlag{
						Name:    "number",
						Usage:   "number of values returned [1, 10000]",
						Aliases: []string{"N"},
						Value:   1,
					},
					&cli.BoolFlag{
						Name:    "unique",
						Usage:   "if true strings are unique, characters may repeat. Returns error if N < (to - from)",
						Aliases: []string{"u"},
					},
				},
				Action: randstr.String(&cfg),
			},
			{
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
				Action: uuid.UUID(&cfg),
			},
			{
				Name:  "blob",
				Usage: "generate random Binary Large OBject. Total size must not exceed 1,048,576 bits (128 Kib)",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "size",
						Usage:   "size of blobs in bits [1, 1048576] must be divisible by 8",
						Aliases: []string{"s"},
						Value:   64,
					},
					&cli.BoolFlag{
						Name:        "hex",
						Usage:       "if true generated data has hex format, base64 otherwise",
						DefaultText: "base64",
					},
					&cli.IntFlag{
						Name:    "number",
						Usage:   "number of values returned [1, 10000]",
						Aliases: []string{"N"},
						Value:   1,
					},
				},
				Action: blob.BLOB(&cfg),
			},
			{
				Name:    "status",
				Usage:   "get current apiKey usage",
				Aliases: []string{"st"},
				Action:  status.Status(&cfg),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
