package randapi

import (
	"encoding/json"
	"fmt"

	"github.com/bohdanch-w/rand-api/entities"
)

func handleErrorResponse(data []byte) (error, error) {
	var errResp entities.ErrorResponse

	if err := json.Unmarshal(data, &errResp); err != nil {
		return nil, err
	}

	return fmt.Errorf("%d - %s", errResp.Error.Code, errResp.Error.Message), nil
}
