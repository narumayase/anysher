# anysher

Una librería de Go que proporciona una forma flexible de crear repositorios para el envío de mensajes. Ofrece dos implementaciones principales: una para enviar mensajes a un tema de Kafka y otra para enviar mensajes a un punto final HTTP. La implementación deseada se elige en función de la configuración proporcionada.

## Características

*   **Interfaz `ProducerRepository`**: Define una interfaz común para el envío de mensajes, lo que permite implementaciones intercambiables.
*   **Implementación de Kafka**: Incluye un `KafkaRepository` que envía mensajes a un tema de Kafka específico.
*   **Implementación HTTP**: Incluye un `HTTPRepository` que envía mensajes a un punto final HTTP configurado.

## Uso

Para utilizar la librería, puedes crear una instancia de `KafkaRepository` o `HTTPRepository`, según tus necesidades.

### Ejemplo: Creación de un Repositorio Kafka

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
	payload := kafka.Payload{
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

### Ejemplo: Creación de un Repositorio HTTP

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