package repository

import (
	"context"

	"github.com/nantapop-kj/go-auth-service/models"
	"gorm.io/gorm"
)

type LoginHistoryRepository interface {
	Create(ctx context.Context, entry *models.UsersLoginHistory) error
}

type LoginHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewLoginHistoryRepository(db *gorm.DB) LoginHistoryRepository {
	return &LoginHistoryRepositoryImpl{db: db}
}

func (r *LoginHistoryRepositoryImpl) Create(ctx context.Context, entry *models.UsersLoginHistory) error {
	return r.db.WithContext(ctx).Create(entry).Error
}
