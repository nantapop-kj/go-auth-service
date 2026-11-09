package bootstrap

import (
	"log"

	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/db"
	"github.com/nantapop-kj/go-auth-service/repository"
	"github.com/nantapop-kj/go-auth-service/service"
)

type Dependencies struct {
	Config      *config.Config
	Tokens      *service.TokenService
	AuthService *service.AuthService
}

func InitDependencies() (*Dependencies, error) {
	cfg := config.Load()

	secret, err := config.LoadJWTSecret(cfg)
	if err != nil {
		return nil, err
	}

	database, err := db.Connect(cfg.Postgres)
	if err != nil {
		return nil, err
	}
	if err := db.Migrate(database); err != nil {
		return nil, err
	}
	log.Println("postgres connected and migrated")

	redisClient, err := config.ConnectRedis(cfg.Redis)
	if err != nil {
		return nil, err
	}
	log.Println("redis connected")

	redisRepo := repository.NewRedisRepo(redisClient)

	tokens := service.NewTokenService(cfg, secret, redisRepo)

	authService := service.NewAuthService(
		cfg,
		repository.NewUsersRepository(database),
		repository.NewLoginHistoryRepository(database),
		redisRepo,
		tokens,
	)

	return &Dependencies{Config: cfg, Tokens: tokens, AuthService: authService}, nil
}
