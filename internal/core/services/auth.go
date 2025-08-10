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

// GenerateTokens hem access hem refresh token üretir
func (as *authService) GenerateTokens(user *domain.User) (string, string, error) {
	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"uuid":    user.UUID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(as.secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"uuid": user.UUID,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(as.secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh token'ı Redis'e kaydet (UUID -> refresh token)
	if err := as.redis.Set(context.Background(), user.UUID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens ile yeni access ve refresh token üret
func (as *authService) RefreshTokens(refreshToken string) (string, string, error) {
	// Token doğrula
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("geçersiz imzalama yöntemi")
		}
		return []byte(as.secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("geçersiz refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["uuid"] == nil {
		return "", "", errors.New("geçersiz claim")
	}
	uuid := claims["uuid"].(string)

	// Redis'ten refresh token'ı doğrula
	storedToken, err := as.redis.Get(context.Background(), uuid)
	if err != nil || storedToken == nil || storedToken.(string) != refreshToken {
		return "", "", errors.New("refresh token bulunamadı veya eşleşmedi")
	}

	// Yeni tokenlar üret
	newAccessClaims := jwt.MapClaims{
		"uuid": uuid,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}
	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims).SignedString([]byte(as.secretKey))
	if err != nil {
		return "", "", err
	}

	newRefreshClaims := jwt.MapClaims{
		"uuid": uuid,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshClaims).SignedString([]byte(as.secretKey))
	if err != nil {
		return "", "", err
	}

	// Redis'te refresh token'ı güncelle
	if err := as.redis.Set(context.Background(), uuid, newRefreshToken); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (as *authService) ValidateToken(token string) (bool, error) {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("geçersiz imzalama yöntemi")
		}
		return []byte(as.secretKey), nil
	})
	return err == nil, err
}

func (as *authService) Logout(ctx context.Context, token string) error {
	return as.redis.Del(ctx, token)
}
