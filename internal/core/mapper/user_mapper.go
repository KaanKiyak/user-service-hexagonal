package mapper

import (
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/dto"
)

func ToLoginResponse(token string) *dto.LoginResponse {
	return &dto.LoginResponse{
		Message:     "Giriş başarılı",
		AccessToken: token,
	}
}
func ToDeleteUserResponse() dto.DeleteUserResponse {
	return dto.DeleteUserResponse{
		Message: "User Deletad Seccesfully",
	}
}
func ToUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
		Role:  user.Role,
	}
}
func ToUserResponseList(users []*domain.User) []dto.UserResponse {
	result := make([]dto.UserResponse, len(users))
	for i, user := range users {
		result[i] = ToUserResponse(user)
	}
	return result
}
func ToUserDomain(dto *dto.CreateUserRequest) *domain.User {
	return &domain.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Age:      dto.Age,
		Password: dto.Password,
		Role:     dto.Role,
	}
}
func ToUodateUserDomain(dto *dto.UpdateUserRequest) *domain.User {
	return &domain.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
