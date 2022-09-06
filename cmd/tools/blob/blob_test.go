package blob_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/bohdanch-w/rand-api/cmd/tools/blob"
	helpers_test "github.com/bohdanch-w/rand-api/cmd/tools/helpers"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services/mock"
)

func TestBlobCommand_SuccessNoParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateBlobs",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateBlobs", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/blob_request_default.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(helpers_test.TestRandResult(t, `["Ox8NH7tk4HM="]`), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{";\x1f\r\x1f\xbbd\xe0s"}, helpers_test.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := blob.NewBlobCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "blob"})
	require.NoError(t, err)
}

func TestBlobCommand_SuccessWithParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateBlobs",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateBlobs", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/blob_request.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(helpers_test.TestRandResult(t, `["c5d13a3b", "6db0943e", "322a961a"]`), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput(
				[]any{"\xc5\xd1:;", "m\xb0\x94>", "2*\x96\x1a"},
				helpers_test.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := blob.NewBlobCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "blob", "-s", "32", "-hex", "-N", "3"})
	require.NoError(t, err)
}

func TestBlobCommand_BadParams(t *testing.T) {
	appConfig := &config.AppConfig{
		APIKey:  "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout: time.Second * 5,
	}

	command := blob.NewBlobCommand(appConfig)
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
			params:        []string{"-s", "-8"},
			expectedError: "`size` param is invalid: must be no less than 1",
		},
		{
			params:        []string{"-s", "1048584"},
			expectedError: "`size` param is invalid: must be no greater than 1048576",
		},
		{
			params:        []string{"-s", "13"},
			expectedError: "`size` parameter must be divisible by 8",
		},
		{
			params:        []string{"-s", "65536", "-N", "1000"},
			expectedError: "size exceed: 65536000 > 1048576",
		},
	}

	for _, tc := range testcases {
		err := app.Run(append([]string{"main.go", "blob"}, tc.params...))
		require.EqualError(t, err, tc.expectedError)
	}
}

func TestBlobCommand_RetrieveFailed(t *testing.T) {
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

	command := blob.NewBlobCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "blob"})
	require.ErrorIs(t, err, entities.Error("test error"))
}

func TestBlobCommandOutputFailed(t *testing.T) {
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
			Return(helpers_test.TestRandResult(t, `["Ox8NH7tk4HM="]`), nil),
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

	command := blob.NewBlobCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "blob"})
	require.ErrorIs(t, err, entities.Error("test error"))
}
