package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/db"
)

func main() {
	database := db.ConnectDB()
	db.AutoMigrate(database)
	app := fiber.New()

	port := config.Getenv("BACKEND_PORT", "3001")
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
