package api

import (
	"database/sql"
	"log"
	"user-service-hexagonal/pkg/auth"

	"github.com/gofiber/fiber/v2"
	"user-service-hexagonal/internal/adapters/handlers"
	"user-service-hexagonal/internal/adapters/repository"
	"user-service-hexagonal/internal/core/services"
)

func main() {
	// 1. DB bağlantısı
	dsn := "root:12345678@tcp(localhost:3306)/user_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB bağlantısı açılamadı: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("DB ping başarısız: %v", err)
	}

	// 2. Redis Adapter'ı oluştur

	// 3. Repository ve Service oluştur
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// 4. Handler'lar
	registerHandler := handlers.NewRegister(userService)
	readUserHandler := handlers.NewReadUser(userService)
	readUsersHandler := handlers.NewReadUsers(userService)
	updateUserHandler := handlers.NewUpdateUser(userService)
	deleteUserHandler := handlers.NewDeleteUser(userService)

	// 5. Fiber app başlat
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		//loglama

	})

	// 6. Route'lar
	api := app.Group("/api")
	api.Post("/users", registerHandler.CreateUser)
	api.Post("/user/login", registerHandler.CreateUser)
	
	user := api.Use(func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token yok",
			})
		}

		// token validation
		valid, err := auth.ValidateToken(token)

		return c.Next()
	})

	user.Get("/users/:id", readUserHandler.ReadUser)
	user.Get("/users", readUsersHandler.ReadUsers)
	user.Put("/users/:id", updateUserHandler.UpdateUser)
	user.Delete("/users/:id", deleteUserHandler.DeleteUser)

	api.Get("/refresh-token", func(c *fiber.Ctx) error {

		auth.RefreshTokens(c)
		return nil

	})

	// 7. Server başlat
	log.Fatal(app.Listen(":8080"))
}
