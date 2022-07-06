package randapi

import (
	"time"

	"github.com/google/uuid"
)

type RandResponse struct {
	ID             uuid.UUID          `json:"id"`
	JsonrpcVersion string             `json:"jsonrpc"`
	Result         RandResponseResult `json:"result"`
}

type RandResponseResult struct {
	Random struct {
		Data      []interface{}
		Timestamp time.Time
	}
	BitsUsed      uint64
	BitsLeft      uint64
	RequestsLeft  uint64
	AdvisoryDelay uint64
}
