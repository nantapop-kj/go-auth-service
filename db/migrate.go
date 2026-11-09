package db

import (
	"fmt"

	"github.com/nantapop-kj/go-auth-service/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return fmt.Errorf("create uuid-ossp extension: %w", err)
	}

	if err := db.AutoMigrate(&models.Users{}, &models.UsersLoginHistory{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}
