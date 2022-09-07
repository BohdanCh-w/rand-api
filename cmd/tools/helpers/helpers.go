package helpers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/google/uuid"
)

// nolint: gomnd
func TestRandResult(t *testing.T, str string) entities.RandResponseResult {
	t.Helper()

	return entities.RandResponseResult{
		Random: entities.RandomData{
			Data:      json.RawMessage([]byte(str)),
			Timestamp: entities.RandTime(time.Date(2022, 8, 25, 12, 15, 44, 395, time.UTC)),
		},
		BitsUsed:      150,
		BitsLeft:      1477,
		RequestsLeft:  233,
		AdvisoryDelay: 1,
	}
}

// nolint: gomnd
func TestRandAPIInfo(t *testing.T, id uuid.UUID) entities.APIInfo {
	t.Helper()

	return entities.APIInfo{
		ID:           id,
		Timestamp:    time.Date(2022, 8, 25, 12, 15, 44, 395, time.UTC),
		BitsUsed:     150,
		BitsLeft:     1477,
		RequestsLeft: 233,
	}
}
