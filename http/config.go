package http

import (
	"github.com/rs/zerolog"
	"strings"
)

// Config contains the application configuration for HTTP.
type Config struct {
	LogLevel string
}

// NewConfiguration creates a new Config instance for HTTP implementation.
// It takes the desired log level as input.
// The provided logLevel string is used to set the global zerolog level.
func NewConfiguration(logLevel string) Config {
	setLogLevel(logLevel)
	return Config{
		LogLevel: logLevel,
	}
}

// setLogLevel sets the global zerolog log level based on the provided string.
// If the logLevel string is not recognized, it defaults to zerolog.InfoLevel.
func setLogLevel(logLevel string) {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
	levelEnv := strings.ToLower(logLevel)

	level, ok := levels[levelEnv]
	if !ok {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}