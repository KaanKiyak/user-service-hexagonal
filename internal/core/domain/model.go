package domain

import "time"

type User struct {
	UUID     string `json:"uuid"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type EventLog struct {
	ID          int       `json:"id"`
	UserID      *int      `json:"user_id,omitempty"`
	Email       string    `json:"email,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	EventType   string    `json:"event_type"` // LOGIN, LOGOUT, PROFILE_REQUEST
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent,omitempty"`
	Status      string    `json:"status"` // SUCCESS, FAILED
	Reason      string    `json:"reason,omitempty"`
	RequestPath string    `json:"request_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
