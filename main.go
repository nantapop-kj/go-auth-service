package main

import "github.com/nantapop-kj/go-auth-service/db"

func main() {
	database := db.ConnectDB()
	db.AutoMigrate(database)
}
