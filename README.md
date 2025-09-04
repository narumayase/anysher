# anysher

A Go library of reusable components for building applications. It includes HTTP, Kafka and Redis clients, Gin-Gonic middlewares, etc.

## Features

*   **HTTP Client**: A wrapper around Go's `net/http` client to simplify making POST HTTP requests.
*   **Kafka Producer**: A client for sending messages to a Kafka topic.
*   **Logging**: A helper to set the global log level for `zerolog`.
*   **Redis**: A client for saving data into Redis.
*   **Gin Middlewares**: A collection of middlewares for the Gin-Gonic framework:
    *   `CORS`: Configures Cross-Origin Resource Sharing.
    *   `Logger`: Logs incoming HTTP requests.
    *   `ErrorHandler`: Handles panics and returns a standardized JSON error response.
    *   `HeadersToContext`: Injects request headers into the `context`.
    *   `RequestIDToLogger`: Adds a request ID to the logger context for better traceability.
    *   `gateway.Sender`: Sends the response to a configured gateway. 

## Usage

The library is divided into packages. You can import the ones you need.

### Example: Using Middlewares in Gin

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/narumayase/anysher/log"
	"github.com/narumayase/anysher/middleware"
	"github.com/narumayase/anysher/middleware/gateway"
)

func main() {
	// Set the global log level
	log.SetLogLevel()

	// Create a new Gin router
	router := gin.New()

	// Use the middlewares
	router.Use(middleware.RequestIDToLogger())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())
	router.Use(gateway.Sender())

	// Define a sample route
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "world"})
	})

	// Start the server
	router.Run(":8080")
}
```

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

### Example: Creating an HTTP Client

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
	cfg := http.NewConfiguration()

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

### Example: Using Redis cache

```go
package main

import (
	"context"
	"github.com/narumayase/anysher/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg := redis.NewConfiguration()

	// Create Redis repository
	redisRepository := redis.NewRedisRepository(cfg)

	// Save data
	if err := redisRepository.Save(context.Background(), "a_key", []byte("{some_data}")); err != nil {
		log.Panic().Msgf("Failed to send data to Redis: %v", err)
	}
	// Retrieve data
	if data, err := redisRepository.Get(context.Background(), "a_key"); err != nil {
		log.Panic().Msgf("Failed to send data to Redis: %v", err)
	} else {
		log.Printf("Successfully retrieve data from Redis %s", data)
	}
}
```