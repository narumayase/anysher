package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

// Producer is an interface that wraps the confluent-kafka-go producer.
type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Events() chan kafka.Event
	Flush(timeoutMs int) int
	Close()
}

type Message struct {
	Key     string
	Headers map[string]string
	Content []byte
}

// Repository Kafka repository.
type Repository struct {
	producer Producer
	topic    string
}

func NewRepository(cfg Config) (*Repository, error) {
	if cfg.KafkaBroker == "" {
		log.Warn().Msg("Kafka broker is not configured; Kafka is disabled.")
		return nil, nil
	}
	p, err := newProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.KafkaBroker})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Info().Msgf("Successfully created Kafka producer for brokers: %s", cfg.KafkaBroker)

	return &Repository{
		producer: p,
		topic:    cfg.KafkaTopic,
	}, nil
}

var newProducer func(configMap *kafka.ConfigMap) (Producer, error) = func(configMap *kafka.ConfigMap) (Producer, error) {
	return kafka.NewProducer(configMap)
}

// Send a message to a Kafka topic.
func (r *Repository) Send(ctx context.Context, payload Message) error {
	if r.producer == nil {
		log.Warn().Msg("Kafka producer is not initialized; cannot send messages.")
		return nil
	}

	var kafkaHeaders []kafka.Header
	for k, v := range payload.Headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key: k, Value: []byte(v),
		})
	}
	log.Debug().Msgf("sending message to Kafka: %v", payload)

	deliveryChan := make(chan kafka.Event)
	err := r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &r.topic, Partition: kafka.PartitionAny},
		Value:          payload.Content,
		Headers:        kafkaHeaders,
		Key:            []byte(payload.Key),
	}, deliveryChan)

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	}
	log.Debug().Msgf("delivered message to topic %s [%d] at offset %v",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	close(deliveryChan)
	return nil
}

// Close closes the Kafka producer.
func (r *Repository) Close() {
	if r.producer != nil {
		r.producer.Close()
	}
}
