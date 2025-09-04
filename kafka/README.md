# kafka

*   **Kafka Producer**: A client for sending messages to a Kafka topic.

## Usage

### Configuration

Create a `.env` file:

- `LOG_LEVEL`: zerolog level.
- `KAFKA_TOPIC`: Kafka topic name to produce.
- `KAFKA_BROKER`: Kafka broker.

### Example: Creating a Kafka Producer

```go
package main

import (
	"context"
	"github.com/narumayase/anysher/kafka"
	"github.com/rs/zerolog/log"
)

func main() {
	// Create Kafka configuration
	cfg := kafka.NewConfiguration()

	// Create a new Kafka repository
	kafkaRepo, err := kafka.NewRepository(cfg)
	if err != nil {
		//log.Fatalf("Failed to create Kafka http: %v", err)
	}
	defer kafkaRepo.Close()

	// Create a payload
	payload := kafka.Message{
		Key:     "somekey",
		Headers: map[string]string{"correlation_id": "123456"},
		Content: []byte("Hello, Kafka!"),
	}

	// Send the message
	if err := kafkaRepo.Send(context.Background(), payload); err != nil {
		log.Err(err).Msg("Failed to send message to Kafka")
	}
}
```