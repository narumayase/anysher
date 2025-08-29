package kafka

import (
	"context"
	"github.com/narumayase/anysher/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaRepository_KafkaDisabled(t *testing.T) {
	cfg := config.Config{}
	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	assert.Nil(t, repo)
}

func TestKafkaRepository_Produce_NilProducer(t *testing.T) {
	repo := &Repository{}
	ctx := context.Background()

	payload := Payload{
		Key:     "key",
		Content: []byte("test message"),
	}

	err := repo.Send(ctx, payload)
	assert.NoError(t, err) // Should return nil error and log a warning
}

func TestKafkaRepository_Close_NilProducer(t *testing.T) {
	repo := &Repository{}
	// Should not panic or cause an error
	repo.Close()
}
