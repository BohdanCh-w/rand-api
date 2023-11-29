package entities_test

import (
	"encoding/json"
	"testing"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/stretchr/testify/require"
)

func TestPregenRandomMarshalling(t *testing.T) {
	testcases := []struct {
		name        string
		in          entities.PregenRand
		expected    string
		expectedErr error
	}{
		{
			name:     "nil",
			expected: `{"a":null}`,
		},
		{
			name: "date",
			in: entities.PregenRand{
				Date: pointerTo("25-05-2022"),
			},
			expected: `{"a":"25-05-2022"}`,
		},
		{
			name: "id",
			in: entities.PregenRand{
				ID: pointerTo("idd-0013"),
			},
			expected: `{"a":"idd-0013"}`,
		},
		{
			name: "both",
			in: entities.PregenRand{
				Date: pointerTo("25-05-2022"),
				ID:   pointerTo("idd-0013"),
			},
			expectedErr: entities.Error("only one of date or id is allowed"),
		},
	}

	type temp struct {
		A entities.PregenRand `json:"a"`
	}

	for _, tc := range testcases {
		v := temp{A: tc.in}

		bb, err := json.Marshal(v)

		if tc.expectedErr == nil {
			require.JSONEq(t, tc.expected, string(bb), tc.name)
			require.NoError(t, err)
		} else {
			require.ErrorIs(t, err, tc.expectedErr)
		}
	}
}

func pointerTo[T any](v T) *T {
	return &v
}
