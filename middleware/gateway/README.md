# gateway Sender

* `gateway.Sender`: Sends the response to a configured gateway.

## Usage

### Configuration

Create a `.env` file:

- `LOG_LEVEL`: zerolog level.
- `GATEWAY_ENABLED`: Defines if the prompt will be sent to the gateway instead of LLM (default:false)
- `GATEWAY_API_URL`: Gateway API URL (optional)
- `GATEWAY_IGNORE_ENDPOINTS`: Endpoints separated by pipe to ignore when sending response to `gateway`. eg:
  `GET:health|POST:send`.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/narumayase/anysher/log"
	"github.com/narumayase/anysher/middleware/gateway"
)

func main() {
	// Set the global log level
	log.SetLogLevel()

	// Create a new Gin router
	router := gin.New()

	// Use the middleware
	router.Use(gateway.Sender())

	// Define a sample route
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "world"})
	})

	// Start the server
	router.Run(":8080")
}
```