package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/dto"
	"user-service-hexagonal/internal/core/mapper"
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
	var req dto.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "geçersiz istek formu",
		})
	}
	//buraya validayionslar gelecek
	userDomain := mapper.ToUserDomain(&req)
	createdUser, err := r.userService.CreateUser(userDomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Kullanıcı oluşturulamadı: " + err.Error(),
		})
	}
	response := mapper.ToUserResponse(createdUser)
	return c.Status(fiber.StatusCreated).JSON(response)
}
