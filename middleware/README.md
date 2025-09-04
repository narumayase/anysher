# middlewares

*   **Gin Middlewares**: A collection of middlewares for the Gin-Gonic framework:
    *   `CORS`: Configures Cross-Origin Resource Sharing.
    *   `Logger`: Logs incoming HTTP requests.
    *   `ErrorHandler`: Handles panics and returns a standardized JSON error response.
    *   `HeadersToContext`: Injects request headers into the `context`.
    *   `RequestIDToLogger`: Adds a request ID to the logger context for better traceability.

## Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/narumayase/anysher/log"
	"github.com/narumayase/anysher/middleware"
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

	// Define a sample route
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "world"})
	})

	// Start the server
	router.Run(":8080")
}
```