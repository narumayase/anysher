package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

type Payload struct {
	Key     string
	Headers map[string]string
	Content []byte
}

// Repository Kafka repository.
type Repository struct {
	producer *kafka.Producer
	topic    string
}

func NewRepository(cfg Config) (*Repository, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.KafkaBroker})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Info().Msgf("Successfully created Kafka producer for brokers: %s", cfg.KafkaBroker)

	return &Repository{
		producer: p,
		topic:    cfg.KafkaTopic,
	}, nil
}

// Send a message to a Kafka topic.
func (r *Repository) Send(ctx context.Context, payload Payload) error {
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
