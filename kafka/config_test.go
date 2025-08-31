package kafka

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		broker      string
		topic       string
		logLevel    string
		expectedCfg Config
	}{
		{
			name:     "Valid configuration",
			broker:   "localhost:9092",
			topic:    "test-topic",
			logLevel: "debug",
			expectedCfg: Config{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test-topic",
				logLevel:    "debug",
			},
		},
		{
			name:     "Invalid log level string is stored",
			broker:   "localhost:9092",
			topic:    "test-topic",
			logLevel: "invalid",
			expectedCfg: Config{
				kafkaBroker: "localhost:9092",
				kafkaTopic:  "test-topic",
				logLevel:    "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfiguration(tt.broker, tt.topic, tt.logLevel)
			assert.Equal(t, tt.expectedCfg, cfg)
		})
	}
}
