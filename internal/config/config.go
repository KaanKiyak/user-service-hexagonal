package config

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

var (
	// JWT secret key, environment variable'dan okunacak veya default atanacak
	JWTSecret string
)

func Init() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		// Eğer env yoksa default değer ata (istersen hataya da çevirebilirsin)
		log.Println("Warning: JWT_SECRET environment variable not set, using default secret")
		JWTSecret = "defaultSecretKey123"
	}
}
func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("DB bağlantısı başarılı")
	return nil
}
