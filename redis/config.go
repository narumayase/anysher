package redis

import (
	"github.com/joho/godotenv"
	anysherlog "github.com/narumayase/anysher/log"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

// Config contains the application configuration for Redis.
type Config struct {
	cacheAddress  string
	cachePassword string
	cacheDatabase int
}

// load loads configuration from environment variables or an .env file
// It takes the configuration from environment variables:
// - CACHE_ADDRESS
// - CACHE_PASSWORD
// - CACHE_DATABASE
// - LOG_LEVEL
func load() Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	config := Config{
		cacheAddress:  getEnv("CACHE_ADDRESS", "localhost:6379"),
		cachePassword: getEnv("CACHE_PASSWORD", ""),
		cacheDatabase: getEnvAsInt("CACHE_DATABASE", 0),
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

// getEnvAsInt retrieves environment variable with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Panic().Err(err).Msg("redis repository: error converting value to int")
		}
		return intValue
	}
	return defaultValue
}
