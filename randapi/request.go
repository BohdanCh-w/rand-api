package randapi

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
)

func (svc *RandomOrgRetriever) NewRequest(
	method string,
	params services.RandParameters,
) (entities.RandomRequest, error) {
	const errParameterInvalid = entities.Error("invalid parameters")

	if params == nil {
		return entities.RandomRequest{}, errParameterInvalid
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
