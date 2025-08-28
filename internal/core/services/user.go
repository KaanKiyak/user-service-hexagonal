package services

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"user-service-hexagonal/internal/config"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
	"user-service-hexagonal/pkg/auth"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(r ports.UserRepository) ports.UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(user *domain.User) (*domain.User, error) {
	// Domain kuralları
	if err := user.ValidateBusinessRules(); err != nil {
		return nil, err
	}

	// UUID üret (yoksa)
	if user.UUID == "" {
		user.UUID = uuid.New().String()
	}

	// Şifreyi hash'le
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)

	// Kaydet
	return s.repo.CreateUser(user)
}

func (s *userService) ReadUser(id int) (*domain.User, error) {
	return s.repo.ReadUser(id)
}

func (s *userService) ReadUsers() ([]*domain.User, error) {
	return s.repo.ReadUsers()
}

func (s *userService) UpdateUser(id, email, password string) error {
	// Yeni şifre geldiyse hash'le
	var hashed string
	if password != "" {
		b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashed = string(b)
	}
	return s.repo.UpdateUser(id, email, hashed)
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) LoginUser(email, password string) (string, error) {
	// Kullanıcıyı çek
	user, err := s.repo.LoginUser(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("kullanıcı bulunamadı")
	}

	// Şifre doğrula (hash vs plain)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("şifre hatalı")
	}

	// JWT üret (aynı secret tüm katmanlarda!)
	token, err := auth.GenerateAccessToken(user, config.JWTSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}
