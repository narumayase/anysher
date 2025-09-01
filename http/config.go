package http

import (
	anysherlog "github.com/narumayase/anysher/log"
)

// Config contains the application configuration for HTTP.
// required environment variable is LOG_LEVEL.
type Config struct {
}

// NewConfiguration creates a new Config instance for HTTP implementation.
// It takes the environment variable LOG_LEVEL
func load() Config {
	anysherlog.SetLogLevel()
	return Config{}
}
