package log

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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
			os.Setenv("LOG_LEVEL", tt.logLevel)
			SetLogLevel()
			assert.Equal(t, tt.expectedLvl, zerolog.GlobalLevel())
		})
	}
}
