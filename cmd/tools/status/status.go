package status

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/config"
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

		usage, err := cfg.RandRetriever.GetUsage(ctx, cfg.APIKey)
		if err != nil {
			return fmt.Errorf("get usage: %w", err)
		}

		if err := cfg.OutputProcessor.GenerateUsageOutput(usage); err != nil {
			return fmt.Errorf("generate usage output: %w", err)
		}

		return nil
	}
}
