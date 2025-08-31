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
// It takes the configuration from environment variables:
// - KAFKA_BROKER
// - KAFKA_TOPIC
// - LOG_LEVEL
func NewConfiguration() Config {
	return load()
}

// load loads configuration from environment variables or an .env file
func load() Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	config := Config{
		kafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		kafkaTopic:  getEnv("KAFKA_TOPIC", "a-topic"),
		logLevel:    getEnv("LOG_LEVEL", "info"),
	}
	anysherlog.SetLogLevel()
	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
