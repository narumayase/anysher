package http

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewConfiguration tests the NewConfiguration function.
func TestNewConfiguration(t *testing.T) {
	logLevel := "debug"

	config := NewConfiguration(logLevel)

	assert.Equal(t, logLevel, config.LogLevel)
}

// TestSetLogLevel tests the setLogLevel function.
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