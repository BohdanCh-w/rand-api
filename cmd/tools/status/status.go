package status

import (
	"context"
	"fmt"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/bohdanch-w/rand-api/randapi"
	"github.com/urfave/cli/v2"
)

func NewStatusCommand(cfg *config.AppConfig) *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "get current apiKey usage",
		Aliases: []string{"st"},
		Action:  status(cfg),
	}
}

func status(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		ctx, cancel := context.WithTimeout(cCtx.Context, cfg.Timeout)
		defer cancel()

		usage, err := randapi.RandAPIUsageExecute(ctx, cfg.APIKey)
		if err != nil {
			return fmt.Errorf("get usage: %w", err)
		}

		output.GenerageStatusOutput(usage)

		return nil
	}
}
