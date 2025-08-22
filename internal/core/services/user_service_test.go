package services

import (
	"testing"
	"user-service-hexagonal/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	user.ID = 1
	return user, nil
}
func (m *mockUserRepository) ReadUser(id int) (*domain.User, error)       { return nil, nil }
func (m *mockUserRepository) ReadUsers() ([]*domain.User, error)          { return nil, nil }
func (m *mockUserRepository) UpdateUser(id, email, password string) error { return nil }
func (m *mockUserRepository) DeleteUser(id string) error                  { return nil }
func (m *mockUserRepository) LoginUser(email string) (*domain.User, error) {
	return nil, nil
}

func TestUserService_CreateUser(t *testing.T) {
	repo := &mockUserRepository{}
	service := NewUserService(repo)

	user := &domain.User{
		Name:     "Jane Doe",
		Email:    "jane@gmail.com",
		Age:      30,
		Password: "pass1234",
		Role:     "user",
	}

	createdUser, err := service.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", createdUser.Name)
	assert.NotEmpty(t, createdUser.UUID)
	assert.NotEqual(t, "pass1234", createdUser.Password) // hashlenmiş olmalı
}
