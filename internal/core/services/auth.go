package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
)

type authService struct {
	redis     ports.RedisPorts
	secretKey string
}

func NewAuthService(redis ports.RedisPorts, secretKey string) ports.AuthService {
	return &authService{
		redis:     redis,
		secretKey: secretKey,
	}
}

func (as *authService) Logout(ctx context.Context, token string) error {
	return as.redis.Del(ctx, token)
}
