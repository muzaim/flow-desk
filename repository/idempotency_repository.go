// repository/idempotency_repository.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IdempotencyData struct {
	ResponseCode int    `json:"response_code"`
	ResponseBody string `json:"response_body"`
}

type IdempotencyRepository interface {
	Get(ctx context.Context, key string, userID uint) (*IdempotencyData, error)
	Set(ctx context.Context, key string, userID uint, data *IdempotencyData, ttl time.Duration) error
}

type idempotencyRepository struct {
	rdb *redis.Client
}

func NewIdempotencyRepository(rdb *redis.Client) IdempotencyRepository {
	return &idempotencyRepository{rdb: rdb}
}

func (r *idempotencyRepository) Get(ctx context.Context, key string, userID uint) (*IdempotencyData, error) {
	redisKey := fmt.Sprintf("idempotency:%d:%s", userID, key)

	val, err := r.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	var data IdempotencyData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *idempotencyRepository) Set(ctx context.Context, key string, userID uint, data *IdempotencyData, ttl time.Duration) error {
	redisKey := fmt.Sprintf("idempotency:%d:%s", userID, key)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, redisKey, jsonData, ttl).Err()
}
