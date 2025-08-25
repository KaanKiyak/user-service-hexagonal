package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // <- ekle
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"user-service-hexagonal/internal/app"
)

var (
	db   *sql.DB
	once sync.Once
)

func getDBConnection() *sql.DB {
	once.Do(func() {
		dsn := "root:12345678@tcp(127.0.0.1:3306)/user_db"
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("DB bağlantısı açılamadı: %v", err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("DB ping başarısız: %v", err)
		}
	})
	return db
}

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func getRedisConnection() *redis.Client {
	redisOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
		})
	})
	return redisClient
}

func main() {

	// DB & Redis
	db := getDBConnection()
	defer db.Close()

	redisClient := getRedisConnection()
	defer redisClient.Close()

	app := app.NewApp(db, redisClient)

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Println("Server kapandı: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server: Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("graceful shutdown error :%v", err)
	}
	log.Println("server: Shutdown Server ...")
}
