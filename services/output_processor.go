package services

import "github.com/bohdanch-w/rand-api/entities"

type OutputProcessor interface {
	GenerateRandOutput(data []interface{}, apiInfo entities.APIInfo) error
	GenerateUsageOutput(status entities.UsageStatus) error
}
