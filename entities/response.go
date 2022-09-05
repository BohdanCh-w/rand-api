package entities

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RandResponse struct {
	ID             uuid.UUID           `json:"id"`
	JsonrpcVersion string              `json:"jsonrpc"`
	Result         *RandResponseResult `json:"result"`
	Error          *ErrorResponse      `json:"error"`
}

type RandResponseResult struct {
	Random        RandomData `json:"random"`
	BitsUsed      uint64     `json:"bitsUsed"`
	BitsLeft      uint64     `json:"bitsLeft"`
	RequestsLeft  uint64     `json:"requestsLeft"`
	AdvisoryDelay uint64     `json:"advisoryDelay"`
}

type RandomData struct {
	Data      json.RawMessage `json:"data"`
	Timestamp RandTime        `json:"completionTime"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
