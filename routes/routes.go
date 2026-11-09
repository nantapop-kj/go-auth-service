package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/bootstrap"
	"github.com/nantapop-kj/go-auth-service/handler"
	"github.com/nantapop-kj/go-auth-service/middleware"
)

func RegisterRoutes(app *fiber.App, deps *bootstrap.Dependencies) {
	cfg := deps.Config
	h := handler.NewAuthHandler(cfg, deps.AuthService, deps.Tokens)

	cookieName := ""
	if cfg.CookieEnabled {
		cookieName = cfg.CookieName
	}

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/api/v1/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	auth.Get("/profile", middleware.RequireAuth(deps.Tokens, cookieName), h.Profile)
}
