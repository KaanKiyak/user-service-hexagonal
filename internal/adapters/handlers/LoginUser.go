package handlers

import (
	"github.com/gofiber/fiber/v2"
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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginUser) LoginUser(c *fiber.Ctx) error {
	var req loginRequest

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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Giriş başarılı",
		"access_token": accessToken,
	})
}
