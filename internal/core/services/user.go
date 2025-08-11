package services

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (u *userService) CreateUser(user *domain.User) (*domain.User, error) {
	// UUID oluştur
	user.UUID = uuid.New().String()

	// Password hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Repo'dan kullanıcı oluşturmayı çağır
	return u.repo.CreateUser(user)
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
