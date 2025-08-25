package middleware

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/pkg/auth"
)

// AuthRequired doğrulama yapan middleware
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token yok",
			})
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  "geçersiz token",
				"detail": err.Error(),
			})
		}

		// userID’yi context’e kaydet
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
