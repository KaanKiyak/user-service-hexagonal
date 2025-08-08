package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	_ "user-service-hexagonal/internal/core/ports"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(addr string, password string, db int) *RedisAdapter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisAdapter{
		client: rdb,
	}
}
func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %q: %v", key, err)
	}

	if err := r.client.Set(ctx, key, data, time.Hour*24*30).Err(); err != nil {
		return fmt.Errorf("failed to set value for key %q: %v", key, err)
	}

	return nil
}
func (r *RedisAdapter) Get(ctx context.Context, key string, value interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss for key %q", key)
	} else if err != nil {
		return fmt.Errorf("failed to get value for key %q: %v", key, err)
	}

	// JSON'dan hedef struct'a parse et
	if err := json.Unmarshal([]byte(data), value); err != nil {
		return fmt.Errorf("failed to unmarshal cache value for key %q: %v", key, err)
	}

	return nil
}
func (r *RedisAdapter) Del(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete value for key %q: %v", key, err)
	}
	return nil
}
