package services

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
)

type authService struct {
	redis ports.RedisPorts
}

func NewAuthService(redis ports.RedisPorts) ports.AuthService {
	return &authService{redis: redis}
}
func (as *authService) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"uuid":    user.UUID,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("Dave"))
}
func (as *authService) ValidateToken(token string) (bool, error) {
	val, err := as.redis.Get(context.Background(), token)
	if err != nil || val == nil {
		return false, err
	}
	return true, nil
}
func (as *authService) Logout(ctx context.Context, token string) error {
	return as.redis.Del(ctx, token)
}
