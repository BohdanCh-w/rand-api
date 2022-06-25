package gausian

import (
	"fmt"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/urfave/cli/v2"
)

func Gausian(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		fmt.Println(cfg)

		return nil
	}
}
