package string

import (
	"fmt"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/urfave/cli/v2"
)

func String(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		fmt.Println(cfg)

		return nil
	}
}