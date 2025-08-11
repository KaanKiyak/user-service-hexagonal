package ports

import (
	"context"
	"user-service-hexagonal/internal/core/domain"
)

type AuthService interface {
	RefreshTokens(refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	Logout(ctx context.Context, token string) error
}
type EventLogRepository interface {
	Save(ctx context.Context, eventLog *domain.EventLog) error
}
type UserService interface {
	CreateUser(user *domain.User) (*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	UpdateUser(id, email, password string) error
	DeleteUser(id string) error
}

type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	UpdateUser(id, email, password string) error
	DeleteUser(id string) error
}
