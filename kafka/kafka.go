package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

// newProducer is a variable that holds the function to create a new Kafka producer.
// This is primarily used for mocking in tests.
var newProducer = func(configMap *kafka.ConfigMap) (Producer, error) {
	return kafka.NewProducer(configMap)
}

// Producer is an interface that wraps the confluent-kafka-go producer.
type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
	Events() chan kafka.Event
	Flush(timeoutMs int) int
	Close()
}

// Message represents the structure of a message to be sent to Kafka.
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

// NewRepository creates a new Kafka repository instance.
// It initializes a Kafka producer if a broker is configured.
// If KafkaBroker is empty, Kafka functionality is disabled and a nil repository is returned.
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

// Send a message to a Kafka topic.
func (r *Repository) Send(ctx context.Context, payload Message) error {
	if r.producer == nil {
		log.Ctx(ctx).Warn().Msg("Kafka producer is not initialized; cannot send messages.")
		return nil
	}

	var kafkaHeaders []kafka.Header
	// Convert message headers to Kafka headers format.
	for k, v := range payload.Headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key: k, Value: []byte(v),
		})
	}
	log.Ctx(ctx).Debug().Msgf("sending message content Kafka: %s", string(payload.Content))
	log.Ctx(ctx).Info().Msgf("sending headers to Kafka: %v", payload.Headers)
	log.Ctx(ctx).Info().Msgf("sending key to Kafka: %s", payload.Key)

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

	// Wait for message delivery report.
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	}
	log.Ctx(ctx).Info().Msgf("delivered message to topic %s [%d] at offset %v",
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
