package handlers

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"user-service-hexagonal/internal/core/mapper"
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
	// 1. ID parametresini al ve int'e çevir
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "geçersiz user id",
		})
	}

	// 2. UserService üzerinden kullanıcıyı çek
	user, err := r.userService.ReadUser(id)
	if err != nil {
		// DB'de kullanıcı yoksa 404 dön
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		// Diğer hatalar için 500 dön
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 3. Mapper ile JSON response oluştur
	return c.JSON(mapper.ToUserResponse(user))
}
