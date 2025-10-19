package db

import (
	"log"

	"github.com/nantapop-kj/go-auth-service/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("❌ Failed to create uuid-ossp extension: %v", err)
	}

	modelsToMigrate := []interface{}{
		&models.Users{},
		&models.UsersLoginHistory{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("❌ [Database] AutoMigrate failed for model %T: %v", model, err)
		}
	}

	log.Println("\033[32m✅ [Database] All migrations completed successfully.\033[0m")
}
