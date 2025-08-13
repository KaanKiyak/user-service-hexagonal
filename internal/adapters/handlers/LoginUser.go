package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/dto"
	"user-service-hexagonal/internal/core/mapper"
	"user-service-hexagonal/internal/core/ports"
)

type LoginUser struct {
	userService ports.UserService
}

func NewLoginUser(userService ports.UserService) *LoginUser {
	return &LoginUser{
		userService: userService,
	}
}

func (l *LoginUser) LoginUser(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// 1️ Body parse
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz istek formatı",
		})
	}

	// 2️ Basit validation
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email ve şifre zorunludur",
		})
	}

	// 3️ Service katmanına login isteği
	accessToken, err := l.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		// Burada err mesajını doğrudan döndürmek yerine genel hata mesajı verilebilir
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Geçersiz email veya şifre",
		})
	}

	// 4️ Başarılı dönüş
	return c.Status(fiber.StatusOK).JSON(mapper.ToLoginResponse(accessToken))

}
