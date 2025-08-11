package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"user-service-hexagonal/internal/adapters/handlers"
	"user-service-hexagonal/internal/adapters/repository"
	"user-service-hexagonal/internal/config"
	"user-service-hexagonal/internal/core/services"
	"user-service-hexagonal/pkg/auth"
	"user-service-hexagonal/pkg/logger"
)

func main() {
	// 1. DB bağlantısı
	dsn := "root:12345678@tcp(localhost:3306)/user_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB bağlantısı açılamadı: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping başarısız: %v", err)
	}

	// 2. Redis Adapter
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	redisAdapter := repository.NewRedisAdapter(redisClient) // ports.RedisPorts implementasyonu
	redisService := services.NewRedisService(redisAdapter)

	// 3. Repository ve Service
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// 4. Handler
	registerHandler := handlers.NewRegister(userService)
	readUserHandler := handlers.NewReadUser(userService)
	readUsersHandler := handlers.NewReadUsers(userService)
	updateUserHandler := handlers.NewUpdateUser(userService)
	deleteUserHandler := handlers.NewDeleteUser(userService)

	// 5. Fiber app başlat
	app := fiber.New()

	// 6. Middleware - Logger
	app.Use(func(c *fiber.Ctx) error {
		var userID *int = nil
		email := ""

		err := logger.LogEvent(
			userID,
			email,
			c.Get("X-Session-ID"),
			"PROFILE_REQUEST",
			"SUCCESS",
			"",
			c.IP(),
			string(c.Request().Header.UserAgent()),
			c.Path(),
		)

		if err != nil {
			// Hata durumunda yapılacaklar
			log.Printf("Logger middleware error: %v", err)
		}

		return c.Next()
	})

	// 7. Routes
	api := app.Group("/api")

	// Register & Login
	api.Post("/users", registerHandler.CreateUser)
	api.Post("/user/login", registerHandler.CreateUser) // login handler yazılmalı

	// Protected routes
	user := api.Use(func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token yok",
			})
		}

		valid, err := auth.ValidateToken(token)
		if err != nil || !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "geçersiz token",
			})
		}

		return c.Next()
	})

	user.Get("/users/:id", readUserHandler.ReadUser)
	user.Get("/users", readUsersHandler.ReadUsers)
	user.Put("/users/:id", updateUserHandler.UpdateUser)
	user.Delete("/users/:id", deleteUserHandler.DeleteUser)

	// Refresh token
	api.Get("/refresh-token", func(c *fiber.Ctx) error {
		refreshToken := c.Query("refresh_token")
		if refreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "refresh token yok",
			})
		}
		secretKey := config.JWTSecret
		newAccess, newRefresh, err := auth.RefreshTokens(refreshToken, secretKey, redisService)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"access_token":  newAccess,
			"refresh_token": newRefresh,
		})
	})

	// 8. Server başlat
	log.Fatal(app.Listen(":8080"))
}
