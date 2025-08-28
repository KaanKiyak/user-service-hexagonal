package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	intapp "user-service-hexagonal/internal/app"
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

<<<<<<< Updated upstream
	app := intapp.NewApp(db, redisClient)

	// create
	body := `{"name":"selo","email":"selo@gmail.com","age":99,"password":"1234","role":"user"}`
=======
	app := app.NewApp(db, redisClient)
	body := `{"name":"neco","email":"neco@gmail.com","age":99,"password":"1234","role":"user"}`
>>>>>>> Stashed changes
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	testUser := &domain.User{ID: 61, Email: "selo@gmail.com", UUID: "e2e-uuid-123"}
	token, _ := auth.GenerateAccessToken(testUser, config.JWTSecret)
	req.Header.Set("Authorization", "Bearer "+token) // public olsa da dursun

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// oluşturulan kullanıcının ID'sini DB'den çekelim
	var createdID int
	err = db.QueryRowContext(context.Background(),
		"SELECT id FROM users WHERE email = ?", "selo@gmail.com").
		Scan(&createdID)
	if err != nil {
		t.Fatalf("created user id alınamadı: %v", err)
	}

	// read (DOĞRU PATH ve DOĞRU ID)
	reqGet := httptest.NewRequest("GET", fmt.Sprintf("/api/users/%d", createdID), nil)
	reqGet.Header.Set("Authorization", "Bearer "+token)
	respGet, err := app.Test(reqGet, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, http.StatusOK, respGet.StatusCode)

	// read ALL
	reqGetAll := httptest.NewRequest("GET", "/api/users", nil)
	reqGetAll.Header.Set("Authorization", "Bearer "+token)
	respGetAll, err := app.Test(reqGetAll, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, http.StatusOK, respGetAll.StatusCode)

	// update
	bodyUpdate := `{"name":"melihat","age":31}`
	reqUpdate := httptest.NewRequest("PUT", fmt.Sprintf("/api/users/%d", createdID), strings.NewReader(bodyUpdate))
	reqUpdate.Header.Set("Content-Type", "application/json")
	reqUpdate.Header.Set("Authorization", "Bearer "+token)

	respUpdate, err := app.Test(reqUpdate, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, http.StatusOK, respUpdate.StatusCode)

	// delete
	reqDel := httptest.NewRequest("DELETE", fmt.Sprintf("/api/users/%d", createdID), nil)
	reqDel.Header.Set("Authorization", "Bearer "+token)

	respDel, err := app.Test(reqDel, -1)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	assert.Equal(t, http.StatusOK, respDel.StatusCode)

	// cleanup (email ile silmek istersen)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM users WHERE email = ?", "selo@gmail.com")
}
