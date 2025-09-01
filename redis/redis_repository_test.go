package redis

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisRepository_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	key := "test:123"
	data := []byte("hola mundo")

	// Mock de Set
	mock.ExpectSet(key, data, 24*time.Hour).SetVal("OK")

	err := repo.Save(ctx, key, data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Mock de Get
	mock.ExpectGet(key).SetVal(string(data))

	result, err := repo.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, string(data), result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisRepository_SaveRedisError(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	data := []byte("hola")
	mock.ExpectSet("key", data, 24*time.Hour).SetErr(redis.ErrClosed)

	err := repo.Save(ctx, "key", data)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisRepository_GetRedisError(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	mock.ExpectGet("key").SetErr(redis.ErrClosed)

	result, err := repo.Get(ctx, "key")
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisRepository_GetKeyNotFound(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	mock.ExpectGet("missing").RedisNil()

	result, err := repo.Get(ctx, "missing")
	assert.ErrorIs(t, err, redis.Nil)
	assert.Equal(t, "", result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRedisRepository(t *testing.T) {
	cfg := Config{
		cacheAddress:  "localhost:6379",
		cachePassword: "",
		cacheDatabase: 0,
	}
	repo := NewRedisRepository(cfg)
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.client)
}
