package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nantapop-kj/go-auth-service/bootstrap"
	"github.com/nantapop-kj/go-auth-service/middleware"
	"github.com/nantapop-kj/go-auth-service/routes"
)

func main() {
	deps, err := bootstrap.InitDependencies()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName:     deps.Config.AppName,
		ProxyHeader: deps.Config.ProxyHeader,
	})
	app.Use(recover.New())
	app.Use(middleware.CORS(deps.Config.CORSOrigins))
	routes.RegisterRoutes(app, deps)

	log.Fatal(app.Listen(":" + deps.Config.ServerPort))
}
