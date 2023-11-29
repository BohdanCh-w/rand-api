package uuid_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	randuuid "github.com/bohdanch-w/rand-api/cmd/tools/uuid"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/pkg/testutils"
	"github.com/bohdanch-w/rand-api/services/mock"
)

func TestUuidCommand_SuccessNoParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateUUIDs",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateUUIDs", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/uuid_request_default.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(testutils.TestRandResult(t, `["03e5a40e-c88d-4b2e-9714-adc581c9437f"]`), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput(
				[]any{uuid.MustParse("03e5a40e-c88d-4b2e-9714-adc581c9437f")},
				testutils.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := randuuid.NewUUIDCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "uuid"})
	require.NoError(t, err)
}

func TestUuidCommand_SuccessWithParams(t *testing.T) { // nolint: funlen
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateUUIDs",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateUUIDs", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/uuid_request.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(testutils.TestRandResult(
				t, `["df428d4f-fd25-4984-ba6c-7557f6b1a853",`+
					` "362097de-5ad9-4f28-ad8c-c5e0ee9c1915",`+
					` "c07a3142-f816-474d-a9a3-0b4512f598ab"]`),
				nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput(
				[]any{
					uuid.MustParse("df428d4f-fd25-4984-ba6c-7557f6b1a853"),
					uuid.MustParse("362097de-5ad9-4f28-ad8c-c5e0ee9c1915"),
					uuid.MustParse("c07a3142-f816-474d-a9a3-0b4512f598ab"),
				},
				testutils.TestRandAPIInfo(t, req.ID)).
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

	command := randuuid.NewUUIDCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "uuid", "-N", "3"})
	require.NoError(t, err)
}

func TestUuidCommand_BadParams(t *testing.T) {
	appConfig := &config.AppConfig{
		APIKey:  "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout: time.Second * 5,
	}

	command := randuuid.NewUUIDCommand(appConfig)
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
	}

	for _, tc := range testcases {
		err := app.Run(append([]string{"main.go", "uuid"}, tc.params...))
		require.EqualError(t, err, tc.expectedError)
	}
}

func TestUuidCommand_RetrieveFailed(t *testing.T) {
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

	command := randuuid.NewUUIDCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "uuid"})
	require.ErrorIs(t, err, entities.Error("test error"))
}

func TestUuidCommandOutputFailed(t *testing.T) {
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
			Return(testutils.TestRandResult(t, `["03e5a40e-c88d-4b2e-9714-adc581c9437f"]`), nil),
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

	command := randuuid.NewUUIDCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "uuid"})
	require.ErrorIs(t, err, entities.Error("test error"))
}
