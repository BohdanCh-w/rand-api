package gausian_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/cmd/tools/gausian"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/pkg/testutils"
	"github.com/bohdanch-w/rand-api/services/mock"
)

func TestGausianCommand_SuccessNoParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateGaussians",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateGaussians", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/gausian_request_default.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(testutils.TestRandResult(t, "[-0.244817]"), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{-0.244817}, testutils.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := gausian.NewGausianCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "gaus"})
	require.NoError(t, err)
}

func TestGausianCommand_SuccessWithParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateGaussians",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateGaussians", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/gausian_request.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(testutils.TestRandResult(t, "[0.224795, 0.582253, -0.580953]"), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{0.224795, 0.582253, -0.580953}, testutils.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
		PregenRand: entities.PregenRand{
			ID: testutils.Pointer("pregen"),
		},
	}

	command := gausian.NewGausianCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "gaus", "-m", "3.14", "-d", "2.46", "-s", "9", "-N", "3"})
	require.NoError(t, err)
}

func TestGausianCommand_BadParams(t *testing.T) {
	appConfig := &config.AppConfig{
		APIKey:  "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout: time.Second * 5,
	}

	command := gausian.NewGausianCommand(appConfig)
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
			params:        []string{"-m", "-1000000.4"},
			expectedError: "`mean` param is invalid: must be no less than -1e+06",
		},
		{
			params:        []string{"-d", "1000000.4"},
			expectedError: "`deviation` param is invalid: must be no greater than 1e+06",
		},
		{
			params:        []string{"-s", "1"},
			expectedError: "`signdig` param is invalid: must be no less than 2",
		},
		{
			params:        []string{"-s", "16"},
			expectedError: "`signdig` param is invalid: must be no greater than 14",
		},
	}

	for _, tc := range testcases {
		err := app.Run(append([]string{"main.go", "gaus"}, tc.params...))
		require.EqualError(t, err, tc.expectedError)
	}
}

func TestGausianCommand_RetrieveFailed(t *testing.T) {
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

	command := gausian.NewGausianCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "gaus"})
	require.ErrorIs(t, err, entities.Error("test error"))
}

func TestGausianCommandOutputFailed(t *testing.T) {
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
			Return(testutils.TestRandResult(t, "[0]"), nil),
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

	command := gausian.NewGausianCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "gaus"})
	require.ErrorIs(t, err, entities.Error("test error"))
}
