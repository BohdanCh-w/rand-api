package version

import (
	"log"

	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/config"
)

func NewVersionCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "get rand-api basic info",
		Action: version(cfg),
	}
}

func version(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		logger := log.Default()
		logger.SetFlags(0)

		logger.Printf("rand-api v%s; go version go1.18.5\n", cfg.Version)

		return nil
	}
}
