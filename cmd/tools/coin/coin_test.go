package coin_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/cmd/tools/coin"
	helpers_test "github.com/bohdanch-w/rand-api/cmd/tools/helpers"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services/mock"
)

func TestCoinCommand_SuccessNoParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateIntegers",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateIntegers", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/coin_request_default.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(helpers_test.TestRandResult(t, "[1]"), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{"tails"}, helpers_test.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := coin.NewCoinCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "coin"})
	require.NoError(t, err)
}

func TestCoinCommand_SuccessWithParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateIntegers",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateIntegers", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/coin_request.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(helpers_test.TestRandResult(t, "[1, 0, 1, 1, 0]"), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{1, 0, 1, 1, 0}, helpers_test.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := coin.NewCoinCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "coin", "-number", "5", "-format", "num"})
	require.NoError(t, err)
}

func TestCoinCommand_BadParams(t *testing.T) {
	appConfig := &config.AppConfig{
		APIKey:  "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout: time.Second * 5,
	}

	command := coin.NewCoinCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	testcases := []struct {
		params        []string
		expectedError string
	}{
		{
			params:        []string{"-number", "-5"},
			expectedError: "`number` param is invalid: must be no less than 1",
		},
		{
			params:        []string{"-N", "10001"},
			expectedError: "`number` param is invalid: must be no greater than 10000",
		},
		{
			params:        []string{"-format", "invalid_value"},
			expectedError: "`format` param is invalid: must be a valid value",
		},
	}

	for _, tc := range testcases {
		err := app.Run(append([]string{"main.go", "coin"}, tc.params...))
		require.EqualError(t, err, tc.expectedError)
	}
}

func TestCoinCommand_RetrieveFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest(gomock.Any(), gomock.Any()).
			Return(entities.RandomRequest{}, entities.Error("test error")),
	)

	appConfig := &config.AppConfig{
		APIKey:        "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:       time.Second * 5,
		RandRetriever: mockRandRetriever,
	}

	command := coin.NewCoinCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "coin"})
	require.ErrorIs(t, err, entities.Error("test error"))
}

func TestCoinCommand_OutputFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest(gomock.Any(), gomock.Any()).
			Return(entities.RandomRequest{}, nil),
		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), gomock.Any()).
			Return(helpers_test.TestRandResult(t, "[0]"), nil),
		mockOutputProcessor.EXPECT().
			GenerateRandOutput(gomock.Any(), gomock.Any()).
			Return(entities.Error("test error")),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := coin.NewCoinCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "coin"})
	require.ErrorIs(t, err, entities.Error("test error"))
}
