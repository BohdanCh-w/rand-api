package services

import "github.com/bohdanch-w/rand-api/entities"

type OutputGenerator interface {
	GenerateRandOutput(data []interface{}, apiInfo entities.APIInfo) error
	GenerateUsageOutput(status entities.UsageStatus) error
}
