package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(token string, secretKey string) (bool, error) {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("geçersiz imzalama yöntemi")
		}
		return []byte(secretKey), nil
	})
	return err == nil, err
}
