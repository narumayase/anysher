package gateway

import "github.com/narumayase/anysher/log"

type Config struct {
	gatewayEnabled bool
	logLevel       string
	gatewayAPIUrl  string
	gatewayToken   string
}

func NewConfig(logLevel string,
	gatewayEnabled bool,
	gatewayAPIUrl string,
	gatewayToken string) *Config {
	log.SetLogLevel(logLevel)
	return &Config{
		gatewayEnabled: gatewayEnabled,
		gatewayAPIUrl:  gatewayAPIUrl,
		gatewayToken:   gatewayToken,
		logLevel:       logLevel,
	}
}
