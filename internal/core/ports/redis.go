package ports

import "context"

type RedisPorts interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key string) error
}
