package e2e

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
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
		log.Printf("sql bağlantı hatası : %v", err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}
func getMockRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
}
func TestUserE2E(t *testing.T) {
	db := getTestDB()
	defer db.Close()

	redisClient := getMockRedis()
	defer redisClient.Close()

	app := app.NewApp(db, redisClient)

	body := `{"name":"neco","email":"neco@gmail.com","age":99,"password":"1234","role":"user"}`
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
	req.Header.Set("Content-type", "application/json")

	testUser := &domain.User{ID: 1, Email: "neco@gmail.com", UUID: "e2e-uuid-123"}
	token, _ := auth.GenerateAccessToken(testUser, config.JWTSecret)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, 201, resp.StatusCode)

	reqGet := httptest.NewRequest("GET", "/api/users/45", nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	respGet, err := app.Test(reqGet, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, http.StatusOK, respGet.StatusCode)
	db.ExecContext(context.Background(), "DELETE FROM users WHERE id = ?", "neco@gmail.com")
}
