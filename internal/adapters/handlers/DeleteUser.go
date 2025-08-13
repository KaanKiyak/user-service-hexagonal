package handlers

import (
	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/core/mapper"
	"user-service-hexagonal/internal/core/ports"
)

type DeleteUser struct {
	userService ports.UserService
}

func NewDeleteUser(userService ports.UserService) *DeleteUser {
	return &DeleteUser{
		userService: userService,
	}
}
func (r *DeleteUser) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is required",
		})
	}
	err := r.userService.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user not found",
		})

	}
	return c.Status(fiber.StatusOK).JSON(mapper.ToDeleteUserResponse())
}
