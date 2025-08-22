package dto

type LoginResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}
type DeleteUserResponse struct {
	Message string `json:"message"`
}
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Age      int    `json:"age" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	Role  string `json:"role"`
}
