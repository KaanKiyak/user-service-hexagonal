package app

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
	"user-service-hexagonal/internal/adapters/handlers"
	"user-service-hexagonal/internal/adapters/repository"
	"user-service-hexagonal/internal/config"
	"user-service-hexagonal/internal/core/services"
	"user-service-hexagonal/pkg/auth"
	"user-service-hexagonal/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

func getDBConnection() *sql.DB {
	once.Do(func() {
		dsn := "root:12345678@tcp(127.0.0.1:3306)/user_db"
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("DB bağlantısı açılamadı: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("DB ping başarısız: %v", err)
		}
	})
	return db
}

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func getRedisConnection() *redis.Client {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	})
	return redisClient
}

func SetupApp() *fiber.App {
	// DB & Redis
	db := getDBConnection()
	redisClient := getRedisConnection()

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

	eventLogger := logger.New(db)
	app := fiber.New()

	// Middleware
	logMiddleware := func(c *fiber.Ctx) error {
		userID := 0
		email := ""
		sessionID := c.Get("X-Session-ID")
		userAgent := string(c.Request().Header.UserAgent())
		path := c.Path()
		reason := ""

		if err := eventLogger.LogEvent(
			userID,
			email,
			sessionID,
			"PROFILE_REQUEST",
			c.IP(),
			userAgent,
			"SUCCESS",
			reason,
			path,
			time.Now(),
		); err != nil {
			log.Printf("Logger middleware error: %v", err)
		}
		return c.Next()
	}

	// Routes
	api := app.Group("/api")

	api.Post("/users", registerHandler.CreateUser)
	api.Post("/user/login", logMiddleware, loginHandler.LoginUser)

	user := api.Use(func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token yok",
			})
		}
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		claims, err := auth.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  "geçersiz token",
				"detail": err.Error(),
			})
		}
		c.Locals("userID", claims.UserID)
		return c.Next()
	})

	user.Get("/users/:id", logMiddleware, readUserHandler.ReadUser)
	user.Get("/users", logMiddleware, readUsersHandler.ReadUsers)
	user.Put("/users/:id", logMiddleware, updateUserHandler.UpdateUser)
	user.Delete("/users/:id", logMiddleware, deleteUserHandler.DeleteUser)

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

	return app
}
