package randapi_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/randapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TestExecutorSuccess(t *testing.T) {
	const (
		apiKey = "6b81b415-80e9-4481-a5f1-58e354742c00" //nolint: gosec
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		request, err := os.ReadFile("./testdata/usage_request.json")
		require.NoError(t, err)

		response, err := os.ReadFile("./testdata/usage_response.json")
		require.NoError(t, err)

		func() {
			id := gjson.GetBytes(body, "id").String()

			_, err = uuid.Parse(id)
			require.NoError(t, err)

			body, err = sjson.SetBytes(body, "id", "00000000-0000-0000-0000-000000000000")
			require.NoError(t, err)

			response, err = sjson.SetBytes(response, "id", id)
			require.NoError(t, err)
		}()

		require.JSONEq(t, string(request), string(body))

		_, err = w.Write(response)
		require.NoError(t, err)
	}))

	defer ts.Close()

	svc := randapi.NewRandomOrgRetriever(
		ts.URL,
		http.DefaultClient,
		false,
	)

	usage, err := svc.GetUsage(context.Background(), apiKey)
	require.NoError(t, err)

	require.Equal(t, apiKey, usage.APIKey)
	require.Equal(t, "running", usage.Status)
	require.Equal(t, time.Date(2022, time.July, 7, 20, 22, 20, 0, time.UTC), time.Time(usage.CreationTime))
	require.Equal(t, uint64(9000), usage.BitsLeft)
	require.Equal(t, uint64(80), usage.RequestsLeft)
	require.Equal(t, uint64(110000), usage.TotalBits)
	require.Equal(t, uint64(700), usage.TotalRequests)
}

func TestExecutor_FailedResoponse(t *testing.T) { // nolint: funlen
	const (
		apiKey = "6b81b415-80e9-4481-a5f1-58e354742c00" // nolint: gosec
	)

	testcases := []struct {
		name          string
		response      func() string
		transformFunc func(string, string) string
		responseCode  int
		expectedError error
	}{
		{
			name: "missmatching id",
			response: func() string {
				response, err := os.ReadFile("./testdata/executor_response.json")
				require.NoError(t, err)

				return string(response)
			},
			transformFunc: func(string, resp string) string { return resp },
			responseCode:  http.StatusOK,
			expectedError: randapi.ErrRequestResponseMissmatch,
		},
		{
			name: "invalid jsonrpc version",
			response: func() string {
				response, err := os.ReadFile("./testdata/executor_response.json")
				require.NoError(t, err)

				return string(response)
			},
			transformFunc: func(req string, resp string) string {
				id := gjson.Get(req, "id").String()

				resp, err := sjson.Set(resp, "id", id)
				require.NoError(t, err)

				resp, err = sjson.Set(resp, "jsonrpc", "3.0")
				require.NoError(t, err)

				return resp
			},
			responseCode:  http.StatusOK,
			expectedError: randapi.ErrUnexpectedJSONRPSVersion,
		},
		{
			name:          "invalid response",
			response:      func() string { return "" },
			transformFunc: func(string, resp string) string { return resp },
			responseCode:  http.StatusOK,
			expectedError: entities.Error("invalid response: decode response: unexpected end of JSON input"),
		},
		{
			name:          "unexpected error code",
			response:      func() string { return "" },
			transformFunc: func(string, resp string) string { return resp },
			responseCode:  http.StatusInternalServerError,
			expectedError: randapi.ErrUnexpectedStatusCode,
		},
		{
			name: "jsonrpc error",
			response: func() string {
				response, err := os.ReadFile("./testdata/usage_response_fail.json")
				require.NoError(t, err)

				return string(response)
			},
			transformFunc: func(req string, resp string) string {
				id := gjson.Get(req, "id").String()

				resp, err := sjson.Set(resp, "id", id)
				require.NoError(t, err)

				return resp
			},
			responseCode:  http.StatusOK,
			expectedError: randapi.ErrErrorInResponse,
		},
	}

	for _, tc := range testcases {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				response := tc.response()
				response = tc.transformFunc(string(body), response)

				w.WriteHeader(tc.responseCode)
				_, err = w.Write([]byte(response))
				require.NoError(t, err)
			}))

			defer ts.Close()

			svc := randapi.NewRandomOrgRetriever(
				ts.URL,
				http.DefaultClient,
				false,
			)

			result, err := svc.ExecuteRequest(context.Background(), &entities.RandomRequest{ID: uuid.New()})
			require.Empty(t, result, tc.name)
			require.True(t, errors.Is(err, tc.expectedError) || err.Error() == tc.expectedError.Error(), tc.name)
		}()
	}
}
