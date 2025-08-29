package kafka

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	broker := "localhost:9092"
	topic := "test-topic"
	logLevel := "debug"

	config := NewConfiguration(broker, topic, logLevel)

	assert.Equal(t, broker, config.KafkaBroker)
	assert.Equal(t, topic, config.KafkaTopic)
	assert.Equal(t, logLevel, config.LogLevel)
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
