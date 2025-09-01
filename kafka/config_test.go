package kafka

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	os.Setenv("KAFKA_TOPIC", "test-topic")
	os.Setenv("KAFKA_BROKER", "localhost:9092")

	expectedConfig := struct {
		name        string
		broker      string
		topic       string
		expectedCfg Config
	}{

		name:   "Valid configuration",
		broker: "localhost:9092",
		topic:  "test-topic",
		expectedCfg: Config{
			kafkaBroker: "localhost:9092",
			kafkaTopic:  "test-topic",
		},
	}
	cfg := NewConfiguration()
	assert.Equal(t, expectedConfig.expectedCfg, cfg)
}
