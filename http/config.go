package http

import (
	"github.com/narumayase/anysher/log"
)

// Config contains the application configuration for HTTP.
type Config struct {
	logLevel string
}

// NewConfiguration creates a new Config instance for HTTP implementation.
// It takes the desired log level as input.
// The provided logLevel string is used to set the global zerolog level.
func NewConfiguration(logLevel string) Config {
	log.SetLogLevel(logLevel)
	return Config{
		logLevel: logLevel,
	}
}
