package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nantapop-kj/go-auth-service/config"
)

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     config.Getenv("CORS_CLIENT_ORIGINS", ""),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	})
}
