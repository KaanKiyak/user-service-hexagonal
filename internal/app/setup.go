package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"user-service-hexagonal/internal/adapters/handlers"
	"user-service-hexagonal/internal/adapters/middleware"
	"user-service-hexagonal/internal/adapters/repository"
	"user-service-hexagonal/internal/core/services"
	"user-service-hexagonal/pkg/auth"
	"user-service-hexagonal/pkg/logger"
)

// NewApp oluşturur ve tüm handler'ları route'lara bağlar
func NewApp(db *sql.DB, redisClient *redis.Client) *fiber.App {
	// Repositories & Services
	redisAdapter := repository.NewRedisAdapter(redisClient)
	redisService := services.NewRedisService(redisAdapter)

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// Handlers
	registerHandler := handlers.NewRegister(userService)
	loginHandler := handlers.NewLoginUser(userService)
	readUserHandler := handlers.NewReadUser(userService)
	readUsersHandler := handlers.NewReadUsers(userService)
	updateUserHandler := handlers.NewUpdateUser(userService)
	deleteUserHandler := handlers.NewDeleteUser(userService)

	// Logger middleware için event logger
	eventLogger := logger.New(db)
	logMiddleware := middleware.LoggerMiddleware(eventLogger)
	authMiddleware := middleware.AuthRequired()

	app := fiber.New()

	// Public routes
	api := app.Group("/api")
	api.Post("/users", registerHandler.CreateUser)
	api.Post("/user/login", logMiddleware, loginHandler.LoginUser)

	// Protected routes
	api.Use(authMiddleware)
	api.Get("/users/:id", logMiddleware, readUserHandler.ReadUser)
	api.Get("/users", logMiddleware, readUsersHandler.ReadUsers)
	api.Put("/users/:id", logMiddleware, updateUserHandler.UpdateUser)
	api.Delete("/users/:id", logMiddleware, deleteUserHandler.DeleteUser)

	// Refresh token route
	api.Get("/refresh-token", func(c *fiber.Ctx) error {
		refreshToken := c.Query("refresh_token")
		if refreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh token yok"})
		}

		secretKey := "" // config.JWTSecret kullanabilirsin
		newAccess, newRefresh, err := auth.RefreshTokens(refreshToken, secretKey, redisService)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"access_token":  newAccess,
			"refresh_token": newRefresh,
		})
	})

	return app
}
