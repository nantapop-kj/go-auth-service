package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/models"
	"gorm.io/gorm"
)

var ErrDuplicate = errors.New("duplicate record")

type UsersRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.Users, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Users, error)
	Create(ctx context.Context, user *models.Users) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID, at time.Time) error
}

type UsersRepositoryImpl struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{db: db}
}

func (r *UsersRepositoryImpl) FindByEmail(ctx context.Context, email string) (*models.Users, error) {
	return r.findOne(ctx, "email = ?", email)
}

func (r *UsersRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Users, error) {
	return r.findOne(ctx, "id = ?", id)
}

func (r *UsersRepositoryImpl) Create(ctx context.Context, user *models.Users) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicate
	}
	return err
}

func (r *UsersRepositoryImpl) UpdateLastLogin(ctx context.Context, id uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).Model(&models.Users{}).Where("id = ?", id).Update("last_login_at", at).Error
}

func (r *UsersRepositoryImpl) findOne(ctx context.Context, query string, arg any) (*models.Users, error) {
	var user models.Users
	err := r.db.WithContext(ctx).Where(query, arg).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
