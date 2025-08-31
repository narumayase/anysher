package kafka

import (
	"github.com/joho/godotenv"
	anysherlog "github.com/narumayase/anysher/log"
	"github.com/rs/zerolog/log"
	"os"
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
	Load()
	anysherlog.SetLogLevel()
	return Config{
		kafkaBroker: KafkaBroker,
		kafkaTopic:  KafkaTopic,
		logLevel:    logLevel,
	}
}

// Load loads configuration from environment variables or an .env file
func Load() Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	config := Config{
		kafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		kafkaTopic:  getEnv("KAFKA_TOPIC", "a-topic"),
		logLevel:    getEnv("LOG_LEVEL", "info"),
	}
	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
