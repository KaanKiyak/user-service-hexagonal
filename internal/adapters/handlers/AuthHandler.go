package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/dto"
	"user-service-hexagonal/internal/core/mapper"
	"user-service-hexagonal/internal/core/ports"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RefreshToken handler'ı: client'tan refresh token alır, yeni access ve refresh token üretir.
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	if err := c.BodyParser(req); err != nil || req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek veya refresh token eksik",
		})
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token yenileme başarısız: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(
		mapper.ToRefreshTokenResponse(newAccessToken, newRefreshToken),
	)
}
