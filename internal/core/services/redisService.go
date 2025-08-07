package services

import (
	"context"
	"user-service-hexagonal/internal/core/ports"
)

type redisService struct {
	redis ports.RedisPorts
}

func NewRedisService(redis ports.RedisPorts) *redisService {
	return &redisService{}
}
func (r *redisService) Set(ctx context.Context, key string, value interface{}) error {
	return r.redis.Set(ctx, key, value)
}
func (r *redisService) Get(ctx context.Context, key string) (interface{}, error) {
	return r.redis.Get(ctx, key)
}
func (r *redisService) Del(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key)
}
