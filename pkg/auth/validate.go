package auth

import (
	"errors"
	"user-service-hexagonal/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims JWT içindeki claim yapısı
type TokenClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateToken token'ı doğrular ve claims döndürür
func ValidateToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("geçersiz token")
	}

	return claims, nil
}
