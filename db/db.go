package db

import (
	"fmt"
	"log"

	"github.com/nantapop-kj/go-auth-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		config.Getenv("POSTGRES_HOST", ""),
		config.Getenv("POSTGRES_USER", ""),
		config.Getenv("POSTGRES_PASSWORD", ""),
		config.Getenv("POSTGRES_DB", ""),
		config.Getenv("POSTGRES_PORT", "5432"),
		config.Getenv("POSTGRES_TIMEZONE", "Asia/Bangkok"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("\033[32m✅ [Database] database connection established successfully.\033[0m")
	return db
}
