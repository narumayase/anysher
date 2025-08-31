# anysher

Una librería de Go con componentes reutilizables para construir aplicaciones. Incluye clientes para HTTP y Kafka, middlewares para Gin-Gonic y utilidades de logging.

## Características

*   **Cliente HTTP**: Un envoltorio sobre el cliente `net/http` de Go para simplificar la realización de peticiones HTTP de tipo POST.
*   **Productor de Kafka**: Un cliente para enviar mensajes a un tema de Kafka.
*   **Logging**: Una utilidad para configurar el nivel de log global para `zerolog`.
*   **Middlewares para Gin**: Una colección de middlewares para el framework Gin-Gonic:
    *   `CORS`: Configura el Intercambio de Recursos de Origen Cruzado (CORS).
    *   `Logger`: Registra las peticiones HTTP entrantes.
    *   `ErrorHandler`: Maneja los `panics` y devuelve una respuesta de error JSON estandarizada.
    *   `HeadersToContext`: Inyecta las cabeceras de la petición en el `context`.
    *   `RequestIDToLogger`: Añade un ID de petición al contexto del logger para una mejor trazabilidad.
    *   `gateway.Sender`: Envía la respuesta a un gateway configurado.

## Uso

La librería está dividida en paquetes. Puedes importar los que necesites.

### Ejemplo: Uso de Middlewares en Gin

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/narumayase/anysher/log"
	"github.com/narumayase/anysher/middleware"
	"github.com/narumayase/anysher/middleware/gateway"
)

func main() {
	// Configurar el nivel de log global
	log.SetLogLevel()

	// Crear un nuevo router de Gin
	router := gin.New()

	cfg := gateway.New()

	// Usar los middlewares
	router.Use(middleware.RequestIDToLogger())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())
	router.Use(gateway.Sender(cfg))

	// Definir una ruta de ejemplo
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "world"})
	})

	// Iniciar el servidor
	router.Run(":8080")
}
```

### Ejemplo: Creación de un Productor de Kafka

```go
package main

import (
	"context"
	"github.com/narumayase/anysher/kafka"
	"github.com/rs/zerolog/log"
)

func main() {
	// Crear la configuración de Kafka
	cfg := kafka.NewConfiguration()

	// Crear un nuevo repositorio de Kafka
	kafkaRepo, err := kafka.NewRepository(cfg)
	if err != nil {
		//log.Fatalf("Failed to create Kafka http: %v", err)
	}
	defer kafkaRepo.Close()

	// Crear un payload
	payload := kafka.Message{
		Key:     "somekey",
		Headers: map[string]string{"correlation_id": "123456"},
		Content: []byte("Hola, Kafka!"),
	}

	// Enviar el mensaje
	if err := kafkaRepo.Send(context.Background(), payload); err != nil {
		log.Err(err).Msg("Failed to send message to Kafka")
	}
}
```

### Ejemplo: Creación de un Cliente HTTP

```go
package main

import (
	"context"
	"github.com/narumayase/anysher/http"
	"log"
	nethttp "net/http"
)

func main() {
	// Crear la configuración HTTP
	cfg := http.NewConfiguration()

	// Crear un nuevo cliente HTTP
	httpClient := http.NewClient(&nethttp.Client{}, cfg)

	// Crear un payload
	payload := http.Payload{
		URL:   "http://localhost:8080",
		Token: "a_bearer_token",
		Headers: map[string]string{"Content-Type": "application/json"},
		Content: []byte("{"Hola, HTTP!"}"),
	}
	// Enviar el payload
	if _, err := httpClient.Post(context.Background(), payload); err != nil {
		log.Printf("Failed to send message via HTTP: %v", err)
	}
}
```