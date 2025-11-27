package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/internal/build"
)

func NewVersionCommand() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "get rand-api basic info",
		Action: version(),
	}
}

func version() cli.ActionFunc {
	return func(cCtx *cli.Context) error {
		apiKey := "unset"

		if build.APIKey != "" {
			masked := len(build.APIKey) - int(float64(len(build.APIKey))*0.25) // nolint: mnd
			apiKey = strings.Repeat("*", masked) + build.APIKey[masked:]
		}

		fmt.Fprintln(os.Stdout, "RandAPI:")
		fmt.Fprintf(os.Stdout, "  version: %s\n", build.Version)
		fmt.Fprintf(os.Stdout, "  go:      %s\n", build.GoVersion)
		fmt.Fprintf(os.Stdout, "  at:      %s\n", build.BuiltAt)
		fmt.Fprintf(os.Stdout, "  api_key: %s\n", apiKey)

		return nil
	}
}
