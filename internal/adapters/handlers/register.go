package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
)

type Register struct {
	userService ports.UserService
}

func NewRegister(userService ports.UserService) *Register {
	return &Register{
		userService: userService,
	}
}
func (r *Register) CreateUser(c *fiber.Ctx) error {
	var req domain.User

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "geçersiz istek formu",
		})
	}
	//buraya validayionslar gelecek
	createdUser, err := r.userService.CreateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı oluşturulamadı: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user_id": createdUser.ID,
		"uuid":    createdUser.UUID,
		"email":   createdUser.Email,
		"name":    createdUser.Name,
		"age":     createdUser.Age,
		"role":    createdUser.Role,
	})
}
