package mocks

import "github.com/confluentinc/confluent-kafka-go/kafka"

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
