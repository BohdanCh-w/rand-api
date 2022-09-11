package string_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/urfave/cli/v2"

	helpers_test "github.com/bohdanch-w/rand-api/cmd/tools/helpers"
	randstr "github.com/bohdanch-w/rand-api/cmd/tools/string"
	"github.com/bohdanch-w/rand-api/config"
	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services/mock"
)

func TestStringCommand_SuccessNoParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateStrings",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateStrings", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/string_request_default.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				chars := gjson.GetBytes(encReq, "characters").String()

				encReq, err = sjson.DeleteBytes(encReq, "characters")
				require.NoError(t, err)

				require.ElementsMatch(t, []rune(chars), []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_")) // nolint: lll
				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(helpers_test.TestRandResult(t, `["a"]`), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{"a"}, helpers_test.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := randstr.NewStringCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "str"})
	require.NoError(t, err)
}

func TestStringCommand_SuccessWithParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := entities.RandomRequest{
		ID:             uuid.MustParse("71d996a7-ff3f-4ba1-84bb-f4cad27eafb6"),
		JsonrpcVersion: "3.0",
		Method:         "generateStrings",
		Params:         json.RawMessage([]byte("encoded params...")),
	}

	mockRandRetriever := mock.NewMockRandRetiever(ctrl)
	mockOutputProcessor := mock.NewMockOutputProcessor(ctrl)

	gomock.InOrder(
		mockRandRetriever.EXPECT().
			NewRequest("generateStrings", gomock.Any()).
			Do(func(_ string, req any) {
				data, err := os.ReadFile("./testdata/string_request.json")
				require.NoError(t, err)

				encReq, err := json.Marshal(req)
				require.NoError(t, err)

				chars := gjson.GetBytes(encReq, "characters").String()

				encReq, err = sjson.DeleteBytes(encReq, "characters")
				require.NoError(t, err)

				require.ElementsMatch(t, []rune(chars), []rune("abcdef"))
				require.JSONEq(t, string(data), string(encReq))
			}).
			Return(req, nil),

		mockRandRetriever.EXPECT().
			ExecuteRequest(gomock.Any(), &req).
			Return(helpers_test.TestRandResult(t, `["feaff", "ccabb", "fdfed"]`), nil),

		mockOutputProcessor.EXPECT().
			GenerateRandOutput([]any{"feaff", "ccabb", "fdfed"}, helpers_test.TestRandAPIInfo(t, req.ID)).
			Return(nil),
	)

	appConfig := &config.AppConfig{
		APIKey:          "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout:         time.Second * 5,
		RandRetriever:   mockRandRetriever,
		OutputProcessor: mockOutputProcessor,
	}

	command := randstr.NewStringCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "str", "-l", "5", "-c", "abcccdefea", "-N", "3", "-u"})
	require.NoError(t, err)
}

func TestStringCommand_BadParams(t *testing.T) {
	appConfig := &config.AppConfig{
		APIKey:  "c6418ada-7874-4907-9367-f43c446686d3",
		Timeout: time.Second * 5,
	}

	command := randstr.NewStringCommand(appConfig)
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
			params:        []string{"-l", "-1"},
			expectedError: "`length` param is invalid: must be no less than 1",
		},
		{
			params:        []string{"-l", "33"},
			expectedError: "`length` param is invalid: must be no greater than 32",
		},
		{
			params:        []string{"-c", ""},
			expectedError: "`charset` param is invalid: length must be no less than 1",
		},
		{
			params:        []string{"-c", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_АБВГДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯабвгдиіїйклмнопрстуфхцчшщьюя0123456789;!@#$%^&*"}, // nolint: lll
			expectedError: "`charset` param is invalid: must be no greater than 128",
		},
		{
			params:        []string{"-c", "abc", "-l", "3", "-N", "100", "-u"},
			expectedError: "`number` of unique requested values is greater than possible with max possible 27",
		},
	}

	for _, tc := range testcases {
		err := app.Run(append([]string{"main.go", "str"}, tc.params...))
		require.EqualError(t, err, tc.expectedError)
	}
}

func TestStringCommand_RetrieveFailed(t *testing.T) {
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

	command := randstr.NewStringCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "str"})
	require.ErrorIs(t, err, entities.Error("test error"))
}

func TestStringCommandOutputFailed(t *testing.T) {
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
			Return(helpers_test.TestRandResult(t, `["0"]`), nil),
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

	command := randstr.NewStringCommand(appConfig)
	app := &cli.App{
		Name:     "test",
		Commands: []*cli.Command{command},
	}

	err := app.Run([]string{"main.go", "str"})
	require.ErrorIs(t, err, entities.Error("test error"))
}
