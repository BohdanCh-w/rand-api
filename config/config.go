package config

import "time"

type AppConfig struct {
	APIKey  string
	Signed  bool
	Timeout time.Duration
	Output  Output
}

type Output struct {
	Filename  *string
	Separator string
	Verbose   bool
	Quite     bool
}
