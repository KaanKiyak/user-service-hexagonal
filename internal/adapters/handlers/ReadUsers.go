package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/ports"
)

type ReadUsers struct {
	userService ports.UserService
}

func NewReadUsers(userService ports.UserService) *ReadUsers {
	return &ReadUsers{
		userService: userService,
	}
}
func (r *ReadUsers) ReadUsers(c *fiber.Ctx) error {
	users, err := r.userService.ReadUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "users not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}
