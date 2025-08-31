package gateway

import (
	"github.com/joho/godotenv"
	anysherlog "github.com/narumayase/anysher/log"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type Config struct {
	gatewayEnabled bool
	logLevel       string
	gatewayAPIUrl  string
	gatewayToken   string
}

// load loads configuration from environment variables or an .env file
// It takes the configuration from environment variables:
// - GATEWAY_API_URL
// - GATEWAY_ENABLED
// - GATEWAY_TOKEN
// - LOG_LEVEL
func load() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	anysherlog.SetLogLevel()
	return &Config{
		gatewayAPIUrl:  getEnv("GATEWAY_API_URL", "http://anyway:9889"),
		gatewayEnabled: getEnvAsBool("GATEWAY_ENABLED", false),
		gatewayToken:   getEnv("GATEWAY_TOKEN", ""),
		logLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}
