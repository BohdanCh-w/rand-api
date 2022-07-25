package config

import (
	"time"

	"github.com/bohdanch-w/rand-api/services"
)

type AppConfig struct {
	APIKey  string
	Timeout time.Duration

	RandRetriever   services.RandRetiever
	OutputProcessor services.OutputProcessor
}
