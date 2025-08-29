package kafka

import (
	"context"
	"errors"
	"github.com/narumayase/anysher/kafka/mocks"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository_KafkaDisabled(t *testing.T) {
	cfg := Config{}
	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	assert.Nil(t, repo)
}

func TestKafkaRepository_Produce_NilProducer(t *testing.T) {
	repo := &Repository{}
	ctx := context.Background()

	payload := Message{
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

func TestNewRepository_Success(t *testing.T) {
	originalNewProducer := newProducer
	defer func() { newProducer = originalNewProducer }()

	mockProducer := &mocks.MockProducer{}
	newProducer = func(cm *kafka.ConfigMap) (Producer, error) {
		return mockProducer, nil
	}

	cfg := Config{KafkaBroker: "localhost:9092"}
	repo, err := NewRepository(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, mockProducer, repo.producer)
}

func TestNewRepository_ProducerError(t *testing.T) {
	originalNewProducer := newProducer
	defer func() { newProducer = originalNewProducer }()

	expectedErr := errors.New("producer error")
	newProducer = func(cm *kafka.ConfigMap) (Producer, error) {
		return nil, expectedErr
	}

	cfg := Config{KafkaBroker: "localhost:9092"}
	repo, err := NewRepository(cfg)

	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), expectedErr.Error())
}

func TestKafkaRepository_Send_Success(t *testing.T) {
	mockProducer := &mocks.MockProducer{
		ProduceFunc: func(msg *kafka.Message, deliveryChan chan kafka.Event) error {
			go func() {
				deliveryChan <- &kafka.Message{TopicPartition: msg.TopicPartition}
			}()
			return nil
		},
	}
	repo := &Repository{producer: mockProducer, topic: "test-topic"}
	ctx := context.Background()

	payload := Message{
		Key:     "key",
		Headers: map[string]string{"hkey": "hvalue"},
		Content: []byte("test message"),
	}

	err := repo.Send(ctx, payload)
	assert.NoError(t, err)
}

func TestKafkaRepository_Send_ProduceError(t *testing.T) {
	mockProducer := &mocks.MockProducer{
		ProduceFunc: func(msg *kafka.Message, deliveryChan chan kafka.Event) error {
			return errors.New("produce error")
		},
	}
	repo := &Repository{producer: mockProducer, topic: "test-topic"}
	ctx := context.Background()

	payload := Message{
		Key:     "key",
		Content: []byte("test message"),
	}

	err := repo.Send(ctx, payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to produce message")
}

func TestKafkaRepository_Send_DeliveryError(t *testing.T) {
	mockProducer := &mocks.MockProducer{
		ProduceFunc: func(msg *kafka.Message, deliveryChan chan kafka.Event) error {
			go func() {
				deliveryChan <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{Error: errors.New("delivery error")},
				}
			}()
			return nil
		},
	}
	repo := &Repository{producer: mockProducer, topic: "test-topic"}
	ctx := context.Background()

	payload := Message{
		Key:     "key",
		Content: []byte("test message"),
	}

	err := repo.Send(ctx, payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delivery failed")
}

func TestKafkaRepository_Close(t *testing.T) {
	closed := false
	mockProducer := &mocks.MockProducer{
		CloseFunc: func() {
			closed = true
		},
	}
	repo := &Repository{producer: mockProducer}
	repo.Close()
	assert.True(t, closed)
}
