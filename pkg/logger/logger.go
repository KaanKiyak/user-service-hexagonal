package logger

import (
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
	"user-service-hexagonal/internal/config"
)

var (
	appLogger *logrus.Logger
	logOnce   sync.Once
)

func GetLogger() *logrus.Logger {
	logOnce.Do(func() {
		appLogger = logrus.New()
		appLogger.SetFormatter(&logrus.JSONFormatter{})
		appLogger.SetLevel(logrus.InfoLevel)
	})
	return appLogger
}

// EventLog structu event_logs tablosunu temsil eder
type EventLog struct {
	ID          int
	UserID      *int
	Email       string
	SessionID   string
	EventType   string // LOGIN, LOGOUT, PROFILE_REQUEST
	IP          string
	UserAgent   string
	Status      string // SUCCESS, FAILED
	Reason      string
	RequestPath string
	CreatedAt   time.Time
}

// LogEvent, verilen event bilgisini validasyonla birlikte DB'ye yazar
func LogEvent(
	userID *int,
	email, sessionID, eventType, status, reason, ip, userAgent, requestPath string,
) error {
	log := GetLogger()
	// Formatlama
	status = strings.ToUpper(strings.TrimSpace(status))
	eventType = strings.ToUpper(strings.TrimSpace(eventType))

	// ENUM kontrolü
	if status != "SUCCESS" && status != "FAILED" {
		log.Printf("WARN: Invalid status '%s', fallback to 'FAILED'", status)
		status = "FAILED"
	}

	// Log kaydı oluştur
	event := EventLog{
		UserID:      userID,
		Email:       email,
		SessionID:   sessionID,
		EventType:   eventType,
		IP:          ip,
		UserAgent:   userAgent,
		Status:      status,
		Reason:      reason,
		RequestPath: requestPath,
		CreatedAt:   time.Now(),
	}

	// SQL sorgusu
	query := `
		INSERT INTO event_logs (user_id, email, session_id, event_type, ip, user_agent, status, reason, request_path, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query,
		event.UserID,
		event.Email,
		event.SessionID,
		event.EventType,
		event.IP,
		event.UserAgent,
		event.Status,
		event.Reason,
		event.RequestPath,
		event.CreatedAt,
	)

	if err != nil {
		log.Printf("ERROR: Event log save error: %v", err)
		return err
	}

	log.Printf("[AUDIT] %s event logged for email=%s status=%s", event.EventType, event.Email, event.Status)
	return nil
}
