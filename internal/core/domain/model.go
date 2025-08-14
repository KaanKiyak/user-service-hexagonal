package domain

type User struct {
	UUID     string `json:"uuid"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
