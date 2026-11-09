package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/models"
	"github.com/nantapop-kj/go-auth-service/repository"
)

const refreshKeyPrefix = "refresh:"

var ErrInvalidToken = errors.New("invalid or expired token")

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenService struct {
	cfg    *config.Config
	secret []byte
	redis  repository.RedisRepository
}

func NewTokenService(cfg *config.Config, secret []byte, redis repository.RedisRepository) *TokenService {
	return &TokenService{cfg: cfg, secret: secret, redis: redis}
}

func (s *TokenService) AccessTTL() time.Duration  { return s.cfg.AccessTokenTTL }
func (s *TokenService) RefreshTTL() time.Duration { return s.cfg.RefreshTokenTTL }

func (s *TokenService) IssueAccessToken(user *models.Users) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    s.cfg.JWTIssuer,
			Subject:   user.ID.String(),
			Audience:  s.cfg.JWTAudience,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *TokenService) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.cfg.JWTIssuer),
		jwt.WithExpirationRequired(),
	}
	if len(s.cfg.JWTAudience) > 0 {
		opts = append(opts, jwt.WithAudience(s.cfg.JWTAudience[0]))
	}

	claims := &AccessClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return s.secret, nil
	}, opts...)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *TokenService) IssueRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(buf)
	if err := s.redis.Set(ctx, refreshKey(token), userID.String(), s.cfg.RefreshTokenTTL); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}
	return token, nil
}

func (s *TokenService) ConsumeRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, ErrInvalidToken
	}
	val, err := s.redis.GetDel(ctx, refreshKey(token))
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, ErrInvalidToken
	}
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}
	return id, nil
}

func refreshKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return refreshKeyPrefix + hex.EncodeToString(sum[:])
}
