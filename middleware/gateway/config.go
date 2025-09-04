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
	gatewayAPIUrl  string
	gatewayToken   string

	ignoreEndpoints []IgnoreEndpoint
}

type IgnoreEndpoint struct {
	Method string
	Path   string
}

// load loads configuration from environment variables or an .env file
// It takes the configuration from environment variables:
// - GATEWAY_API_URL
// - GATEWAY_ENABLED
// - GATEWAY_TOKEN
// - GATEWAY_IGNORE_ENDPOINTS -> format eg: GET:health|POST:send
// - LOG_LEVEL
func load() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	anysherlog.SetLogLevel()

	return &Config{
		gatewayAPIUrl:   getEnv("GATEWAY_API_URL", "http://anyway:9889"),
		gatewayEnabled:  getEnvAsBool("GATEWAY_ENABLED", false),
		gatewayToken:    getEnv("GATEWAY_TOKEN", ""),
		ignoreEndpoints: getIgnoreEndpoints(),
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

func getIgnoreEndpoints() []IgnoreEndpoint {
	ignore := getEnv("GATEWAY_IGNORE_ENDPOINTS", "")
	var ignoreList []IgnoreEndpoint
	if ignore != "" {
		items := strings.Split(ignore, "|")
		for _, item := range items {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) != 2 {
				log.Printf("Invalid ignore endpoint format: %s", item)
				continue
			}
			ignoreList = append(ignoreList, IgnoreEndpoint{
				Method: strings.ToUpper(parts[0]),
				Path:   parts[1],
			})
		}
	}
	return ignoreList
}
