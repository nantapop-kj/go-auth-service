package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/models"
	"github.com/nantapop-kj/go-auth-service/utils"
	"gorm.io/gorm"
)

type UsersRepository interface {
	FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user *models.Users) (uuid.UUID, error)
}

type UsersRepositoryImpl struct {
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{DB: db}
}

func (r *UsersRepositoryImpl) FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	db, err := utils.ResolveDB(tx, r.DB)
	if err != nil {
		return false, err
	}

	var count int64
	if err := db.WithContext(ctx).Model(&models.Users{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UsersRepositoryImpl) CreateUser(ctx context.Context, tx *gorm.DB, user *models.Users) (uuid.UUID, error) {
	db, err := utils.ResolveDB(tx, r.DB)
	if err != nil {
		return uuid.Nil, err
	}

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}
