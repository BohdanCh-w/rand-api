package entities

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RandomRequest struct {
	ID             uuid.UUID       `json:"id"`
	JsonrpcVersion string          `json:"jsonrpc"`
	Method         string          `json:"method"`
	Params         json.RawMessage `json:"params"`
}
