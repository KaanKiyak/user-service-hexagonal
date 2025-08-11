package services

import (
	"context"
	"user-service-hexagonal/internal/core/ports"
)

type RedisService struct {
	redis ports.RedisPorts
}

func NewRedisService(redis ports.RedisPorts) *RedisService {
	return &RedisService{redis: redis}
}

func (r *RedisService) Set(ctx context.Context, key string, value interface{}) error {
	return r.redis.Set(ctx, key, value)
}

func (r *RedisService) Get(ctx context.Context, key string, value interface{}) error {
	return r.redis.Get(ctx, key, value)
}

func (r *RedisService) Del(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key)
}
