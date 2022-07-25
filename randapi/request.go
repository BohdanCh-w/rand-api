package randapi

import (
	"encoding/json"
	"fmt"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
	"github.com/google/uuid"
)

func (svc *RandomOrgRetriever) NewRequest(
	method string,
	params services.RandParameters,
) (entities.RandomRequest, error) {
	if params == nil {
		return entities.RandomRequest{}, fmt.Errorf("invalid params")
	}

	bb, err := json.Marshal(params)
	if err != nil {
		return entities.RandomRequest{}, fmt.Errorf("marhsal params: %w", err)
	}

	return entities.RandomRequest{
		ID:             uuid.New(),
		JsonrpcVersion: jsonRPCVersion,
		Method:         method,
		Params:         bb,
	}, nil
}
