package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"user-service-hexagonal/internal/core/domain"
)

// RedisClient dışarıdan sağlanacak
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
}

// GenerateTokens hem access hem refresh token üretir
func GenerateTokens(user *domain.User, secretKey string, redis RedisClient) (string, string, error) {
	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"uuid":    user.UUID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"uuid": user.UUID,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh token'ı Redis'e kaydet
	if err := redis.Set(context.Background(), user.UUID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens ile yeni access ve refresh token üret
func RefreshTokens(refreshToken string, secretKey string, redis RedisClient) (string, string, error) {
	// Token doğrula
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("geçersiz imzalama yöntemi")
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("geçersiz refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["uuid"] == nil {
		return "", "", errors.New("geçersiz claim")
	}
	uuid := claims["uuid"].(string)

	// Redis'ten refresh token doğrula
	var storedToken string
	err = redis.Get(context.Background(), uuid, &storedToken)
	if err != nil || storedToken != refreshToken {
		return "", "", errors.New("refresh token bulunamadı veya eşleşmedi")
	}

	// Yeni tokenlar üret
	newAccessClaims := jwt.MapClaims{
		"uuid": uuid,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}
	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	newRefreshClaims := jwt.MapClaims{
		"uuid": uuid,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	// Redis'te refresh token güncelle
	if err := redis.Set(context.Background(), uuid, newRefreshToken); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
