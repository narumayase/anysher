package kafka

import (
	"github.com/rs/zerolog"
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
				KafkaBroker: "localhost:9092",
				KafkaTopic:  "test-topic",
				LogLevel:    "debug",
			},
		},
		{
			name:     "Invalid log level string is stored",
			broker:   "localhost:9092",
			topic:    "test-topic",
			logLevel: "invalid",
			expectedCfg: Config{
				KafkaBroker: "localhost:9092",
				KafkaTopic:  "test-topic",
				LogLevel:    "invalid",
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

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		logLevel    string
		expectedLvl zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"fatal", zerolog.FatalLevel},
		{"panic", zerolog.PanicLevel},
		{"invalid", zerolog.InfoLevel},
		{"", zerolog.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			setLogLevel(tt.logLevel)
			assert.Equal(t, tt.expectedLvl, zerolog.GlobalLevel())
		})
	}
}