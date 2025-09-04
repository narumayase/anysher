# redis

* **Redis**: A client for saving data into Redis.

## Usage

### Configuration

Create a `.env` file:

- `LOG_LEVEL`: zerolog level.
- `CACHE_ADDRESS`: redis address (default:localhost:6379).
- `CACHE_PASSWORD`: redis password.
- `CACHE_DATABASE`: redis database (default:0).

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