# http

*   **HTTP Client**: A wrapper around Go's `net/http` client to simplify making POST HTTP requests.

## Usage

### Configuration

Create a `.env` file:

- `LOG_LEVEL`: zerolog level.

### Example: Creating a HTTP Client

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