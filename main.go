package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/bootstrap"
	"github.com/nantapop-kj/go-auth-service/config"
)

func main() {
	_, err := bootstrap.InitDependencies()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	port := config.Getenv("BACKEND_PORT", "3001")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
