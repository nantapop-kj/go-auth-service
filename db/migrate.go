package db

import (
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	log.Println("\033[32m✅ [Database] Migration completed successfully.\033[0m")
}
