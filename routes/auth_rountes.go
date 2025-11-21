package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/bootstrap"
	"github.com/nantapop-kj/go-auth-service/handler"
	"github.com/nantapop-kj/go-auth-service/repository"
	"github.com/nantapop-kj/go-auth-service/service"
)

func RegisterAuthRoutes(
	group fiber.Router,
	prefix string,
	deps *bootstrap.Dependencies,
) {
	authGroup := group.Group(prefix)
	userRepo := repository.NewUsersRepository(deps.DB)
	svc := service.NewAuthService(deps.DB, deps.Minio, deps.Redis, userRepo)
	authHandler := handler.NewAuthHandler(svc)

	authGroup.Post("/register", authHandler.HandlerRegister)
}
