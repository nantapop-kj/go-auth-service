package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/bootstrap"
)

func RegisterRoutes(app *fiber.App, deps *bootstrap.Dependencies) {
	apiGroup := app.Group("/api/v1")

	RegisterAuthRoutes(apiGroup, "auth", deps)
}
