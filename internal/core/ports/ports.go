package ports

import (
	"context"
	"user-service-hexagonal/internal/core/domain"
)

type AuthService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(token string) (bool, error)
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
	UpdateMembershipStatus(id string, status bool) error
}

type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	ReadUser(id string) (*domain.User, error)
	ReadUsers() ([]*domain.User, error)
	UpdateUser(id, email, password string) error
	DeleteUser(id string) error
	UpdateMembershipStatus(id string, status bool) error
}
