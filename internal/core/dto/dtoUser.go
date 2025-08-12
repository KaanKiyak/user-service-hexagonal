package dto

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type UserResponse struct {
	UUID  string `json:"uuid"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	Role  string `json:"role"`
}
