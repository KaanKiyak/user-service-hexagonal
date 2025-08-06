package ports

import (
	"context"
	"user-service-hexagonal/internal/core/domain"
)

type AuthService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (bool, error)
}
type EventLogRepository interface {
	Save(ctx context.Context, eventLog *domain.EventLog) error
}
type UserService interface {
	RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
type UserRepository interface {
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
}
