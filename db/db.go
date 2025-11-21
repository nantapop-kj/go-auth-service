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
		config.GetEnv("POSTGRES_HOST", ""),
		config.GetEnv("POSTGRES_USER", ""),
		config.GetEnv("POSTGRES_PASSWORD", ""),
		config.GetEnv("POSTGRES_DB", ""),
		config.GetEnv("POSTGRES_PORT", "5432"),
		config.GetEnv("POSTGRES_TIMEZONE", "Asia/Bangkok"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("\033[32m✅ [Database] database connection established successfully.\033[0m")
	return db
}
