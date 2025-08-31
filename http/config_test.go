package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewConfiguration tests the NewConfiguration function.
func TestNewConfiguration(t *testing.T) {
	logLevel := "debug"

	config := NewConfiguration(logLevel)

	assert.Equal(t, logLevel, config.logLevel)
}
