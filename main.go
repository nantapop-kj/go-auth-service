package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/bootstrap"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/middleware"
	"github.com/nantapop-kj/go-auth-service/routes"
)

func main() {
	deps, err := bootstrap.InitDependencies()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(middleware.CORSMiddleware())
	routes.RegisterRoutes(app, deps)

	port := config.GetEnv("BACKEND_PORT", "3001")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
