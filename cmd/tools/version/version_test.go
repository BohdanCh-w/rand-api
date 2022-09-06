package version_test

import (
	"testing"

	"github.com/bohdanch-w/rand-api/cmd/tools/version"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestStatusCommandSuccess(t *testing.T) {
	appConfig := &config.AppConfig{
		Version: "1.0.0",
	}

	command := version.NewVersionCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "version"})
	require.NoError(t, err)
}
