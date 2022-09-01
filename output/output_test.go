package output_test

import (
	"strings"
	"testing"
	"time"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/output"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandOutput(t *testing.T) {
	apiInfo := entities.APIInfo{
		ID:           uuid.MustParse("1b13d22f-d633-4ba4-b667-89f1310a00f6"),
		Timestamp:    time.Date(2022, 8, 25, 12, 0, 0, 0, time.UTC),
		BitsUsed:     350,
		BitsLeft:     249_650,
		RequestsLeft: 999,
	}

	testcases := []struct {
		name           string
		data           []any
		apiInfo        entities.APIInfo
		opVerbose      bool
		opQuiet        bool
		opSeparator    string
		expectedOutput string
	}{
		{
			name:           "default int",
			data:           []any{1, 2, 3},
			apiInfo:        apiInfo,
			opSeparator:    ", ",
			expectedOutput: `1, 2, 3`,
		},
		{
			name:           "default string",
			data:           []any{"abc", "def"},
			apiInfo:        apiInfo,
			opSeparator:    ", ",
			expectedOutput: `abc, def`,
		},
		{
			name:           "separator",
			data:           []any{1, 2, 3},
			apiInfo:        apiInfo,
			opSeparator:    "- | -",
			expectedOutput: `1- | -2- | -3`,
		},
	}

	for _, tc := range testcases {
		rr := &Recorder{}

		outputer := output.NewOutputProcessor(tc.opVerbose, tc.opVerbose, tc.opSeparator, rr)

		err := outputer.GenerateRandOutput(tc.data, tc.apiInfo)
		require.NoError(t, err)

		require.Equal(t, strings.TrimSpace(tc.expectedOutput), strings.TrimSpace(rr.String()))
	}
}
