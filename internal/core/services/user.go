package services

import (
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) CreateUser(email, password string) (*domain.User, error) {
	// Burada hashleme, email validasyonu gibi business logic eklenebilir
	return u.repo.CreateUser(email, password)
}

func (u *userService) ReadUser(id string) (*domain.User, error) {
	return u.repo.ReadUser(id)
}

func (u *userService) ReadUsers() ([]*domain.User, error) {
	return u.repo.ReadUsers()
}

func (u *userService) UpdateUser(id, email, password string) error {
	// Şifre değiştiriliyorsa hashlenebilir burada
	return u.repo.UpdateUser(id, email, password)
}

func (u *userService) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}

func (u *userService) UpdateMembershipStatus(id string, status bool) error {
	return u.repo.UpdateMembershipStatus(id, status)
}
