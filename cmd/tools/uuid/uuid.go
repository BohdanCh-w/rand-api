package uuid

import (
	"fmt"

	"github.com/bohdanch-w/rand-api/config"
	"github.com/urfave/cli/v2"
)

func UUID(cfg *config.AppConfig) cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		fmt.Println(cfg)

		return nil
	}
}
