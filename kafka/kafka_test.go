package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

// MockProducer is a mock implementation of the Producer interface.
type MockProducer struct {
	ProduceFunc func(msg *kafka.Message, deliveryChan chan kafka.Event) error
	EventsFunc  func() chan kafka.Event
	FlushFunc   func(timeoutMs int) int
	CloseFunc   func()
}

func (m *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	if m.ProduceFunc != nil {
		return m.ProduceFunc(msg, deliveryChan)
	}
	return nil
}

func (m *MockProducer) Events() chan kafka.Event {
	if m.EventsFunc != nil {
		return m.EventsFunc()
	}
	return nil
}

func (m *MockProducer) Flush(timeoutMs int) int {
	if m.FlushFunc != nil {
		return m.FlushFunc(timeoutMs)
	}
	return 0
}

func (m *MockProducer) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
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

	mockProducer := &MockProducer{}
	newProducer = func(cm *kafka.ConfigMap) (Producer, error) {
		return mockProducer, nil
	}

	load()
	repo, err := NewRepository()

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

	load()
	repo, err := NewRepository()

	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), expectedErr.Error())
}

func TestKafkaRepository_Send_Success(t *testing.T) {
	mockProducer := &MockProducer{
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
	mockProducer := &MockProducer{
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
	mockProducer := &MockProducer{
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
	mockProducer := &MockProducer{
		CloseFunc: func() {
			closed = true
		},
	}
	repo := &Repository{producer: mockProducer}
	repo.Close()
	assert.True(t, closed)
}

func TestKafkaRepository_Send_NoHeaders(t *testing.T) {
	mockProducer := &MockProducer{
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
		Content: []byte("test message without headers"),
	}

	err := repo.Send(ctx, payload)
	assert.NoError(t, err)
	// Optionally, you could add assertions here to check if headers were indeed empty in the produced message
}
