package randapi

import (
	"encoding/json"

	"github.com/google/uuid"
)

func NewRandomRequest(method string, params json.RawMessage) RandomRequest {
	return RandomRequest{
		ID:             uuid.New(),
		JsonrpcVersion: "2.0",
		Method:         method,
		Params:         params,
	}
}

type RandomRequest struct {
	ID             uuid.UUID       `json:"id"`
	JsonrpcVersion string          `json:"jsonrpc"`
	Method         string          `json:"method"`
	Params         json.RawMessage `json:"params"`
}
