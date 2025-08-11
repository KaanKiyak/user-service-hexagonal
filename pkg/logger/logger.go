package logger

import (
	"database/sql"
	"time"
)

// EventType türü (enum karşılığı)
type EventType string

const (
	EventLogin          EventType = "LOGIN"
	EventLogout         EventType = "LOGOUT"
	EventProfileRequest EventType = "PROFILE_REQUEST"
)

// EventStatus türü (enum karşılığı)
type EventStatus string

const (
	StatusSuccess EventStatus = "SUCCESS"
	StatusFailed  EventStatus = "FAILED"
)

// EventLog kayıt modeli
type EventLog struct {
	UserID      *int
	Email       *string
	SessionID   *string
	EventType   EventType
	IP          string
	UserAgent   *string
	Status      EventStatus
	Reason      *string
	RequestPath *string
	CreatedAt   time.Time
}

type Logger struct {
	db *sql.DB
}

// NewLogger MySQL bağlantısını başlatır
func NewLogger(dsn string) (*Logger, error) {
	// Örnek DSN: "user:password@tcp(localhost:3306)/user_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Bağlantı test et
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Logger{db: db}, nil
}

// LogEvent event_logs tablosuna kayıt ekler
func (l *Logger) LogEvent(log EventLog) error {
	query := `
	INSERT INTO event_logs (
		user_id, email, session_id, event_type, ip, user_agent, status, reason, request_path, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := l.db.Exec(query,
		log.UserID,
		log.Email,
		log.SessionID,
		string(log.EventType),
		log.IP,
		log.UserAgent,
		string(log.Status),
		log.Reason,
		log.RequestPath,
		time.Now(),
	)
	return err
}

// Close MySQL bağlantısını kapatır
func (l *Logger) Close() error {
	return l.db.Close()
}
