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
	RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetProfile(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error)
}
type UserRepository interface {
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error)
}
