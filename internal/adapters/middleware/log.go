package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
	"user-service-hexagonal/pkg/logger"
)

// LoggerMiddleware HTTP isteklerini loglayan middleware
func LoggerMiddleware(eventLogger *logger.EventLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := 0
		email := ""
		sessionID := c.Get("X-Session-ID")
		userAgent := string(c.Request().Header.UserAgent())
		path := c.Path()
		reason := ""

		if err := eventLogger.LogEvent(
			userID,
			email,
			sessionID,
			"PROFILE_REQUEST",
			c.IP(),
			userAgent,
			"SUCCESS",
			reason,
			path,
			time.Now(),
		); err != nil {
			log.Printf("Logger middleware error: %v", err)
		}

		return c.Next()
	}
}
