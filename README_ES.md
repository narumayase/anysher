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
	"anysher/config"
	"anysher/internal/domain"
	"anysher/internal/infrastructure/repository"
	"context"
	"github.com/rs/zerolog/log"
)

func main() {
	// Crear configuración de Kafka
	cfg := config.NewKafkaConfiguration("localhost:9092","un-tema", "info")

	// Crear un nuevo repositorio de Kafka
	kafkaRepo, err := repository.NewKafkaRepository(cfg)
	if err != nil {
		//log.Fatalf("Error al crear el repositorio de Kafka: %v", err)
	}
	defer kafkaRepo.Close()

	// Crear un payload
	payload := domain.Payload{
		KafkaPayload: domain.KafkaPayload{Key: "una-clave"},
		Headers:      map[string]string{"correlation_id": "123456"},
		Content:      []byte("¡Hola, Kafka!"),
	}

	// Enviar el mensaje
	if err := kafkaRepo.Send(context.Background(), payload); err != nil {
		log.Err(err).Msg("Error al enviar el mensaje a Kafka")
	}
}
```

### Ejemplo: Creación de un Repositorio HTTP

```go
package main

import (
	"anysher/internal/domain"
	"anysher/internal/infrastructure/repository"
	"context"
	"log"
	"net/http"
)

func main() {
	// Crear configuración HTTP
	cfg := config.NewHTTPConfiguration("info")
	
	// Crear un nuevo cliente HTTP
	httpClient := repository.NewHttpClient(&http.Client{}, cfg)

	// Crear un payload
	payload := domain.Payload{
		HTTPPayload:  domain.HTTPPayload{
			URL:   "http://localhost:8080",
			Token: "un_bearer_token",
		},
		Headers:      map[string]string{"Content-Type": "application/json"},
		Content:      []byte("{\"¡Hola, HTTP!\"}"),
	}

	// Enviar el mensaje
	if _, err := httpClient.Post(context.Background(), payload); err != nil {
		log.Printf("Error al enviar el mensaje a través de HTTP: %v", err)
	}
}
```