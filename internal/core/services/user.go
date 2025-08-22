package services

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
	"user-service-hexagonal/pkg/auth"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{
		repo: repo,
	}
}

// dto struct (crate user) -> doğrulama fonksiyonları
// dto struct (create user) -> apapter
// dto struct (user login response) -> handler

func (u *userService) CreateUser(user *domain.User) (*domain.User, error) {
	// UUID oluştur
	user.UUID = uuid.New().String()

	if err := user.ValidateBusinessRules(); err != nil {
		return nil, err
	}

	// Password hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Repo'dan kullanıcı oluşturmayı çağır
	return u.repo.CreateUser(user)
}
func (u *userService) ReadUser(id int) (*domain.User, error) {
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
func (u *userService) LoginUser(email, password string) (string, error) {
	user, err := u.repo.LoginUser(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("kullanıcı bulunamadı")
	}

	// Şifre kontrolü
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("şifre hatalı")
	}

	// Sadece access token üret
	return auth.GenerateAccessToken(user, "SECRET_KEY")
}
