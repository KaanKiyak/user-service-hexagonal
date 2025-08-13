package handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"user-service-hexagonal/internal/core/dto"
	"user-service-hexagonal/internal/core/ports"
)

type UpdateUser struct {
	userService ports.UserService
}

func NewUpdateUser(userService ports.UserService) *UpdateUser {
	return &UpdateUser{
		userService: userService,
	}
}
func (r *UpdateUser) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is required",
		})
	}
	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "geçersiz istek",
		})
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "şifre işlenemedi",
			})
			req.Password = string(hashedPassword)
		}
		err = r.userService.UpdateUser(id, req.Email, req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "kullanıcı güncellenemdi" + err.Error()})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"massege": "kullanıcı başarıyla gğncellenemedi",
	})
}
