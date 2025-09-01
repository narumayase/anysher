package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisRepository_SaveMock(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	key := "test:123"
	data := map[string]string{"foo": "bar"}
	jsonData, _ := json.Marshal(data)

	mock.ExpectSet(key, jsonData, 24*time.Hour).SetVal("OK")

	err := repo.Save(ctx, key, data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisRepository_GetMock(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	key := "test:123"
	data := map[string]string{"foo": "bar"}
	jsonData, _ := json.Marshal(data)

	mock.ExpectGet(key).SetVal(string(jsonData))

	result, err := repo.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisRepository_SaveRedisError(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	data := map[string]string{"foo": "bar"}
	jsonData, _ := json.Marshal(data)

	mock.ExpectSet("key", jsonData, 24*time.Hour).SetErr(redis.ErrClosed)

	err := repo.Save(ctx, "key", data)
	assert.Error(t, err)
}

func TestRedisRepository_GetKeyNotFound(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	mock.ExpectGet("missing").RedisNil()

	result, err := repo.Get(ctx, "missing")
	assert.ErrorIs(t, err, redis.Nil)
	assert.Nil(t, result)
}

func TestRedisRepository_GetRedisError(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	mock.ExpectGet("key").SetErr(redis.ErrClosed)

	result, err := repo.Get(ctx, "key")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRedisRepository_GetUnmarshalError(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	repo := &RedisRepository{client: db}

	mock.ExpectGet("key").SetVal("not-json")

	result, err := repo.Get(ctx, "key")
	assert.Error(t, err)
	assert.Nil(t, result)
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
