package config

import (
	"database/sql"
)

var DB *sql.DB

var (
	// JWT secret key, environment variable'dan okunacak veya default atanacak
	JWTSecret = "JWT_SECRET"
)
