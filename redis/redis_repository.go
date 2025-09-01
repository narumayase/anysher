package cache

import (
	"context"
	"encoding/json"
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

// Save saves the metadata into redis
func (r *RedisRepository) Save(ctx context.Context, key string, metadata map[string]string) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to marshal metadata")
		return err
	}
	// TODO ver de hacer configurable esta duración
	if err = r.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to save metadata")
		return err
	}
	log.Ctx(ctx).Debug().Msgf("metadata saved in Redis: %s:%v", key, metadata)
	return nil
}

// Get gets the metadata from redis
func (r *RedisRepository) Get(ctx context.Context, key string) (map[string]string, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch metadata")
		return nil, err
	}
	var metadata map[string]string
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal metadata")
		return nil, err
	}
	log.Ctx(ctx).Debug().Msgf("metadata retrieve from Redis: %s:%v", key, metadata)
	return metadata, nil
}
