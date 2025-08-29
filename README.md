# anysher

A Go library that provides a flexible way to create repositories for sending messages. It offers two main implementations: one for sending messages to a Kafka topic and another for sending messages to an HTTP endpoint. The desired implementation is chosen based on the provided configuration.

## Features

*   **`ProducerRepository` Interface**: Defines a common interface for sending messages, allowing for interchangeable implementations.
*   **Kafka Implementation**: Includes a `KafkaRepository` that sends messages to a specified Kafka topic.
*   **HTTP Implementation**: Includes an `HTTPRepository` that sends messages to a specified HTTP endpoint.

## Usage

To use the library, you can create an instance of either `KafkaRepository` or `HTTPRepository`, depending on your needs.

### Example: Creating a Kafka Repository

```go
package main

import (
	"context"
	"github.com/narumayase/anysher/kafka"
	"github.com/rs/zerolog/log"
)

func main() {
	// Create Kafka configuration
	cfg := kafka.NewConfiguration("localhost:9092", "a-topic", "info")

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

### Example: Creating an HTTP Repository

```go
package main

import (
	"context"
	"github.com/narumayase/anysher/http"
	"log"
	nethttp "net/http"
)

func main() {
	// Create HTTP configuration
	cfg := http.NewConfiguration("info")

	// Create a new HTTP client
	httpClient := http.NewClient(&nethttp.Client{}, cfg)

	// Create a payload
	payload := http.Payload{
		URL:   "http://localhost:8080",
		Token: "a_bearer_token",
		Headers: map[string]string{"Content-Type": "application/json"},
		Content: []byte("{\"Hello, HTTP!\"}"),
	}
	// Post the payload
	if _, err := httpClient.Post(context.Background(), payload); err != nil {
		log.Printf("Failed to send message via HTTP: %v", err)
	}
}
```