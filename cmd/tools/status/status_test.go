package status_test

import (
	"testing"
	"time"

	"github.com/bohdanch-w/rand-api/cmd/tools/status"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestStatusCommandSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usageStatus := entities.UsageStatus{
		APIKey:        "c6418ada-7874-4907-9367-f43c446686d3",
		Status:        "running",
		CreationTime:  entities.RandTime(time.Date(2022, 8, 25, 12, 15, 44, 395, time.UTC)),
		TotalRequests: 250,
		TotalBits:     30000,
		RequestsLeft:  50,
		BitsLeft:      12000,
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			GetUsage(gomock.Any(), "c6418ada-7874-4907-9367-f43c446686d3").
			Return(usageStatus, nil),
		mockOutputProcessor.EXPECT().
			GenerateUsageOutput(usageStatus).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := status.NewStatusCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "status"})
	require.NoError(t, err)
}

func TestStatusCommandRetrieveFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			GetUsage(gomock.Any(), "c6418ada-7874-4907-9367-f43c446686d3").
			Return(entities.UsageStatus{}, entities.Error("test error")),
	)

	appConfig := &config.AppConfig{
		APIKey:        "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:       time.Second * 5,
		RandRetriever: mockRandRetriever,
	}

	command := status.NewStatusCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "status"})
	require.ErrorIs(t, err, entities.Error("test error"))
}

func TestStatusCommandOutputFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			GetUsage(gomock.Any(), "c6418ada-7874-4907-9367-f43c446686d3").
			Return(entities.UsageStatus{}, nil),
		mockOutputProcessor.EXPECT().
			GenerateUsageOutput(entities.UsageStatus{}).
			Return(entities.Error("test error")),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := status.NewStatusCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "status"})
	require.ErrorIs(t, err, entities.Error("test error"))
}
