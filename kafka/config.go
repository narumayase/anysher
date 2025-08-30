package kafka

import (
	"github.com/rs/zerolog"
	"strings"
)

// Config contains the application configuration for Kafka.
type Config struct {
	KafkaBroker string
	KafkaTopic  string
	LogLevel    string
}

// NewConfiguration creates a new Config instance for Kafka implementation.
// It takes Kafka broker address, topic, and desired log level as input.
// The provided logLevel string is used to set the global zerolog level.
func NewConfiguration(KafkaBroker string, KafkaTopic string, logLevel string) Config {
	setLogLevel(logLevel)
	return Config{
		KafkaBroker: KafkaBroker,
		KafkaTopic:  KafkaTopic,
		LogLevel:    logLevel,
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