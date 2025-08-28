package benchmark

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"

	"github.com/redis/go-redis/v9"
	"user-service-hexagonal/internal/core/services" // 🔑 burası eksikti

	"net/http/httptest"
	"strings"
	"testing"
	"user-service-hexagonal/internal/adapters/handlers"
	"user-service-hexagonal/internal/app"
	"user-service-hexagonal/internal/core/domain"
)

type mockRepo struct{}

func (m *mockRepo) CreateUser(user *domain.User) (*domain.User, error) {
	return user, nil
}

func (m *mockRepo) ReadUser(id int) (*domain.User, error) {
	return &domain.User{ID: id, Email: "mock@gmail.com", Name: "Mock"}, nil
}

func (m *mockRepo) ReadUsers() ([]*domain.User, error) {
	users := []*domain.User{}
	for i := 0; i < 100; i++ {
		users = append(users, &domain.User{
			ID:    i,
			Email: fmt.Sprintf("mock%d@gmail.com", i),
			Name:  "Mock",
		})
	}
	return users, nil
}

func (m *mockRepo) UpdateUser(id, email, password string) error {
	return nil
}

func (m *mockRepo) DeleteUser(id string) error {
	return nil
}

func (m *mockRepo) LoginUser(email string) (*domain.User, error) {
	return &domain.User{Email: email, Password: "1234"}, nil
}
func getTestDB() *sql.DB {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/user_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
func getMockRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	return redisClient
}
func BenchmarkRegisterHandler(b *testing.B) {
	db := getTestDB()
	defer db.Close()

	redisClient := getMockRedis()
	defer redisClient.Close()

	app := app.NewApp(db, redisClient)

	body := `{"name":"lana","email":"lana@gmail.com","age":99,"password":"1234","role":"user"}`
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
	req.Header.Set("Content-type", "application/json")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Test(req, -1)
	}

}

// Tek user getiren handler benchmark
func BenchmarkReadUserHandler(b *testing.B) {
	service := services.NewUserService(&mockRepo{})
	readUserHandler := handlers.NewReadUser(service)

	app := fiber.New()
	app.Get("/users/:id", readUserHandler.ReadUser)

	req := httptest.NewRequest("GET", "/users/1", nil)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req, -1)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected status %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

// Tüm user’ları getiren handler benchmark
func BenchmarkReadUsersHandler(b *testing.B) {
	service := services.NewUserService(&mockRepo{})
	readUsersHandler := handlers.NewReadUsers(service)
	app := fiber.New()
	app.Get("/users", readUsersHandler.ReadUsers)
	req := httptest.NewRequest("GET", "/users", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req, -1)
		if err != nil {
			b.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected status %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}
