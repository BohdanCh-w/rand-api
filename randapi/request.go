package randapi

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const jsonRPCVersion = "2.0"

func NewRandomRequest(method string, params interface{}) (RandomRequest, error) {
	if params == nil {
		return RandomRequest{}, fmt.Errorf("invalid params")
	}

	bb, err := json.Marshal(params)
	if err != nil {
		return RandomRequest{}, fmt.Errorf("marhsal params: %w", err)
	}

	return RandomRequest{
		ID:             uuid.New(),
		JsonrpcVersion: jsonRPCVersion,
		Method:         method,
		Params:         bb,
	}, nil
}

type RandomRequest struct {
	ID             uuid.UUID       `json:"id"`
	JsonrpcVersion string          `json:"jsonrpc"`
	Method         string          `json:"method"`
	Params         json.RawMessage `json:"params"`
}
