package services

import (
	"context"

	"github.com/bohdanch-w/rand-api/entities"
)

type RandParameters interface{}

type RandRetiever interface {
	NewRequest(method string, params RandParameters) (entities.RandomRequest, error)
	ExecuteRequest(ctx context.Context, randReq *entities.RandomRequest) (entities.RandResponseResult, error)
	GetUsage(ctx context.Context, apiKey string) (entities.UsageStatus, error)
}
