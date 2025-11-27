package config

import (
	"time"

	"github.com/bohdanch-w/rand-api/entities"
	"github.com/bohdanch-w/rand-api/services"
)

type AppConfig struct {
	APIKey     string
	PregenRand entities.PregenRand
	Timeout    time.Duration

	RandRetriever   services.RandRetiever
	OutputProcessor services.OutputGenerator
}
