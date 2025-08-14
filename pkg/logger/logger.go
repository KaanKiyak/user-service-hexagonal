package logger

import (
	"database/sql"
	"errors"
	"time"
)

type EventLogger struct {
	db *sql.DB
}

func New(db *sql.DB) *EventLogger {
	if db == nil {
		panic("logger.New: db is nil")
	}
	return &EventLogger{db: db}
}

func (l *EventLogger) LogEvent(
	user_id int,
	email string,
	session_id string,
	event_type string,
	ip string,
	user_agent string,
	status string,
	reason string,
	request_path string,
	created_at time.Time,
) error {
	if l.db == nil {
		return errors.New("logger: db is nil")
	}

	// Örnek kolon adları: event_logs tablonu ile eşleşecek şekilde güncelle
	const q = `
		INSERT INTO event_logs
			(user_id, email, session_id, event_type, ip, user_agent, status, reason, request_path, created_at)
		VALUES (?,?,?,?,?,?,?,?,?,?)
	`
	_, err := l.db.Exec(q,
		user_id, email, session_id, event_type, ip, user_agent, status, reason, request_path, created_at,
	)
	return err
}
