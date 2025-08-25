package integration

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"user-service-hexagonal/internal/app"
	"user-service-hexagonal/internal/config"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/pkg/auth"
)

func getTestDB() *sql.DB {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/user_db?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("fb bağlantısı açılamadı: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping başarısız: %v", err)
	}
	return db
}

func getTestRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
}

func TestCreateUserIntegration(t *testing.T) {
	testUser := &domain.User{
		ID:    1,
		Email: "jane@gmail.com",
		UUID:  "test-uuid-123",
	}
	token, err := auth.GenerateAccessToken(testUser, config.JWTSecret)
	db := getTestDB()
	if err := db.Ping(); err != nil {
		t.Fatalf("db ping başarısız : %v", err)
	}
	defer db.Close()
	// Context oluştur
	ctx := context.Background()

	// Redis ping testi

	redisClient := getTestRedis()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		t.Fatalf("Redis ping başarısız: %v", err)
	}

	defer redisClient.Close()
	app := app.NewApp(db, redisClient)
	body := `{"name":"kyle","email":"kyle@gmail.com","age":25,"password":"1234","role":"user"}`
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request error : %v", err)
	}
	assert.Equal(t, 201, resp.StatusCode)
}
