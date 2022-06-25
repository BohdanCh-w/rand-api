package config

import "time"

type AppConfig struct {
	APIKey     string
	Signed     bool
	Timeout    time.Duration
	OutputFile *string
}
