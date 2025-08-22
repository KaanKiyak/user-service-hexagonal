package integration

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
	"user-service-hexagonal/internal/app"
)

func Test_CreateUser_Integration(t *testing.T) {
	app := app.SetupApp()
	body := `{"name":"Jane","email":"jane@gmail.com","age":25,"password":"1234","role":"user"}`
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 201, resp.StatusCode)
}
