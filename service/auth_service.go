package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/dto"
	"github.com/nantapop-kj/go-auth-service/enums"
	"github.com/nantapop-kj/go-auth-service/models"
	"github.com/nantapop-kj/go-auth-service/repository"
	"golang.org/x/crypto/bcrypt"
)

const loginFailKeyPrefix = "login_fail:"

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTooManyAttempts    = errors.New("too many failed login attempts, try again later")
	ErrUserInactive       = errors.New("user is inactive")
	ErrUserNotFound       = errors.New("user not found")
)

type ClientInfo struct {
	IP        string
	UserAgent string
}

type AuthService struct {
	cfg       *config.Config
	users     repository.UsersRepository
	history   repository.LoginHistoryRepository
	redis     repository.RedisRepository
	tokens    *TokenService
	dummyHash []byte
}

func NewAuthService(
	cfg *config.Config,
	users repository.UsersRepository,
	history repository.LoginHistoryRepository,
	redis repository.RedisRepository,
	tokens *TokenService,
) *AuthService {
	dummyHash, _ := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), bcrypt.DefaultCost)
	return &AuthService{cfg: cfg, users: users, history: history, redis: redis, tokens: tokens, dummyHash: dummyHash}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest, client ClientInfo) (*dto.TokenResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.Users{
		Email:        normalizeEmail(req.Email),
		PasswordHash: string(hash),
		IsActive:     true,
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	s.record(ctx, &user.ID, user.Email, enums.LoginActionRegister, client)
	return s.issueTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest, client ClientInfo) (*dto.TokenResponse, error) {
	email := normalizeEmail(req.Email)

	if s.isLocked(ctx, email) {
		s.record(ctx, nil, email, enums.LoginActionLocked, client)
		return nil, ErrTooManyAttempts
	}

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if user == nil {
		_ = bcrypt.CompareHashAndPassword(s.dummyHash, []byte(req.Password))
		return nil, s.loginFailed(ctx, nil, email, client)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, s.loginFailed(ctx, &user.ID, email, client)
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	_ = s.redis.Delete(ctx, loginFailKeyPrefix+email)
	now := time.Now()
	if err := s.users.UpdateLastLogin(ctx, user.ID, now); err != nil {
		log.Printf("update last_login_at: %v", err)
	}
	s.record(ctx, &user.ID, email, enums.LoginActionLogin, client)

	return s.issueTokens(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	userID, err := s.tokens.ConsumeRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	user, err := s.GetUser(ctx, userID)
	if errors.Is(err, ErrUserNotFound) {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string, client ClientInfo) error {
	userID, err := s.tokens.ConsumeRefreshToken(ctx, refreshToken)
	if errors.Is(err, ErrInvalidToken) {
		return nil
	}
	if err != nil {
		return err
	}
	s.record(ctx, &userID, "", enums.LoginActionLogout, client)
	return nil
}

func (s *AuthService) GetUser(ctx context.Context, id uuid.UUID) (*models.Users, error) {
	user, err := s.users.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	return user, nil
}

func (s *AuthService) isLocked(ctx context.Context, email string) bool {
	if s.cfg.LoginMaxAttempts <= 0 {
		return false
	}
	failures, err := s.redis.GetInt(ctx, loginFailKeyPrefix+email)
	return err == nil && failures >= s.cfg.LoginMaxAttempts
}

func (s *AuthService) loginFailed(ctx context.Context, userID *uuid.UUID, email string, client ClientInfo) error {
	if s.cfg.LoginMaxAttempts > 0 {
		if err := s.redis.IncrWithTTL(ctx, loginFailKeyPrefix+email, s.cfg.LoginLockDuration); err != nil {
			log.Printf("count failed login: %v", err)
		}
	}
	s.record(ctx, userID, email, enums.LoginActionFail, client)
	return ErrInvalidCredentials
}

func (s *AuthService) issueTokens(ctx context.Context, user *models.Users) (*dto.TokenResponse, error) {
	access, err := s.tokens.IssueAccessToken(user)
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.IssueRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		AccessToken:      access,
		TokenType:        "Bearer",
		ExpiresIn:        int64(s.tokens.AccessTTL().Seconds()),
		RefreshToken:     refresh,
		RefreshExpiresIn: int64(s.tokens.RefreshTTL().Seconds()),
		User:             dto.NewUserResponse(user),
	}, nil
}

func (s *AuthService) record(ctx context.Context, userID *uuid.UUID, email string, action enums.LoginAction, client ClientInfo) {
	entry := &models.UsersLoginHistory{
		UserID:    userID,
		Email:     email,
		IPAddress: client.IP,
		UserAgent: truncate(client.UserAgent, 255),
		Action:    action,
	}
	if err := s.history.Create(ctx, entry); err != nil {
		log.Printf("write login history: %v", err)
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}
