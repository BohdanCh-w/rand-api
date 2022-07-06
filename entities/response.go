package entities

import (
	"encoding/json"
	"fmt"
	"strings"
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
		Data      json.RawMessage `json:"data"`
		Timestamp completionTime  `json:"completionTime"`
	} `json:"random"`
	BitsUsed      uint64 `json:"bitsUsed"`
	BitsLeft      uint64 `json:"bitsLeft"`
	RequestsLeft  uint64 `json:"requestsLeft"`
	AdvisoryDelay uint64 `json:"advisoryDelay"`
}

type completionTime time.Time

func (c *completionTime) UnmarshalJSON(data []byte) error {
	const format = "2006-01-02 15:04:05Z"

	str := strings.Trim(strings.TrimSpace(string(data)), "\"")

	t, err := time.Parse(format, str)
	if err != nil {
		return fmt.Errorf("parse completionTime: %w", err)
	}

	*c = completionTime(t)

	return nil
}
