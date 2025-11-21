package utils

import (
	"errors"

	"gorm.io/gorm"
)

func ResolveDB(tx *gorm.DB, defaultDB *gorm.DB) (*gorm.DB, error) {
	db := tx
	if db == nil {
		db = defaultDB
	}

	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	return db, nil
}
