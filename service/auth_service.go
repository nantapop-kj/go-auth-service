package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/dto"
	"github.com/nantapop-kj/go-auth-service/models"
	"github.com/nantapop-kj/go-auth-service/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB              *gorm.DB
	Minio           repository.MinioRepository
	Redis           repository.RedisRepository
	JwtSecret       []byte
	UsersRepository repository.UsersRepository
}

func NewAuthService(
	DB *gorm.DB,
	minio repository.MinioRepository,
	redis repository.RedisRepository,
	usersRepository repository.UsersRepository,
) *AuthService {
	return &AuthService{
		DB:              DB,
		Minio:           minio,
		Redis:           redis,
		JwtSecret:       []byte(config.GetEnv("JWT_SECRET_KEY", "")),
		UsersRepository: usersRepository,
	}
}

func (s *AuthService) ServiceRegister(ctx context.Context, req dto.RegisterRequest) (uuid.UUID, string, error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		return uuid.Nil, "", tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	rollback := func(err error) (uuid.UUID, string, error) {
		tx.Rollback()
		return uuid.Nil, "", err
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.ProfileImage == nil {
		return rollback(errors.New("profile image is required"))
	}

	if req.ProfileImage.Size > 10*1024*1024 {
		return rollback(errors.New("profile image file too large"))
	}

	existing, err := s.UsersRepository.FindUserByEmail(ctx, tx, req.Email)
	if err != nil {
		return rollback(err)
	}

	if existing {
		return rollback(errors.New("email already registered"))
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return rollback(errors.New("failed to hash password"))
	}

	fileContent, err := req.ProfileImage.Open()
	if err != nil {
		return rollback(fmt.Errorf("failed to read profile image: %w", err))
	}
	defer fileContent.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(fileContent)
	if err != nil {
		return rollback(fmt.Errorf("failed to read file content: %w", err))
	}

	url, err := s.Minio.UploadImage(ctx, req.ProfileImage.Filename, buf.Bytes())
	if err != nil {
		return rollback(fmt.Errorf("failed to upload avatar: %w", err))
	}

	// url, err := s.Minio.GetFileURL(ctx, key)
	// if err != nil {
	// 	return rollback(fmt.Errorf("failed to generate avatar URL: %w", err))
	// }

	user := &models.Users{
		Email:           req.Email,
		Password:        string(hashed),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		ProfileImageURL: url,
	}

	fmt.Println("A100 user", user)

	userID, err := s.UsersRepository.CreateUser(ctx, tx, user)
	if err != nil {
		return rollback(errors.New("failed to create user"))
	}

	token, err := config.GenerateJWT(userID, user.Email)
	if err != nil {
		return rollback(err)
	}

	if err := tx.Commit().Error; err != nil {
		return uuid.Nil, "", err
	}

	return userID, token, nil
}
