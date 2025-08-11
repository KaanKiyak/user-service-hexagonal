package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"user-service-hexagonal/internal/core/ports"
)

type redisAdapter struct {
	client *redis.Client
}

// NewRedisAdapter Redis adapter constructor
func NewRedisAdapter(client *redis.Client) ports.RedisPorts {
	return &redisAdapter{client: client}
}

// Set key-value (stringify)
func (r *redisAdapter) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

// Get value by key
func (r *redisAdapter) Get(ctx context.Context, key string, value interface{}) error {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	// Burada interface{} parametreyi pointer bekleyip yazıyoruz
	switch v := value.(type) {
	case *string:
		*v = result
	default:
		// Eğer tip string değilse sadece raw olarak döndür
		// Tip dönüşümü burada senin kullanımına göre genişletilebilir
		return nil
	}
	return nil
}

// Delete key
func (r *redisAdapter) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
