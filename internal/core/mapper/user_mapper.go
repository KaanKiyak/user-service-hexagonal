package mapper

import (
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/dto"
)

// DTO -> domain
func ToUserDomain(dto *dto.CreateUserRequest) *domain.User {
	return &domain.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Age:      dto.Age,
		Password: dto.Password,
		Role:     dto.Role,
	}
}
func ToUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
		Role:  user.Role,
	}
}
