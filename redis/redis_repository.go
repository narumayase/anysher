package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

// RedisRepository implements the CacheRepository interface using Redis.
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new instance of RedisRepository.
func NewRedisRepository(config Config) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.cacheAddress,
		Password: config.cachePassword,
		DB:       config.cacheDatabase,
	})
	return &RedisRepository{
		client: rdb,
	}
}

// Save saves the data into redis
func (r *RedisRepository) Save(ctx context.Context, key string, data []byte) error {
	// TODO ver de hacer configurable esta duración
	if err := r.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to save metadata")
		return err
	}
	log.Ctx(ctx).Debug().Msgf("metadata saved in Redis: %s:%v", key, string(data))
	return nil
}

// Get gets the data from redis
func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch metadata")
		return "", err
	}
	log.Ctx(ctx).Debug().Msgf("metadata retrieve from Redis: %s:%v", key, data)
	return data, nil
}
