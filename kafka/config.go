package kafka

import (
	"github.com/narumayase/anysher/log"
)

// Config contains the application configuration for Kafka.
type Config struct {
	kafkaBroker string
	kafkaTopic  string
	logLevel    string
}

// NewConfiguration creates a new Config instance for Kafka implementation.
// It takes Kafka broker address, topic, and desired log level as input.
// The provided logLevel string is used to set the global zerolog level.
func NewConfiguration(KafkaBroker string, KafkaTopic string, logLevel string) Config {
	log.SetLogLevel(logLevel)
	return Config{
		kafkaBroker: KafkaBroker,
		kafkaTopic:  KafkaTopic,
		logLevel:    logLevel,
	}
}
