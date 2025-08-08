package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/ports"
)

type ReadUser struct {
	userService ports.UserService
}

func NewReadUser(userService ports.UserService) *ReadUser {
	return &ReadUser{
		userService: userService,
	}
}
func (r *ReadUser) ReadUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is required",
		})
	}
	user, err := r.userService.ReadUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	return c.JSON(user)
}
