package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LucasAndFlores/user_api/config"
	"github.com/LucasAndFlores/user_api/database"
	"github.com/LucasAndFlores/user_api/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var PORT = fmt.Sprintf(":%v", os.Getenv("PORT"))

func main() {
	err := config.LoadEnvVariables()

	if err != nil {
		log.Fatalf("An error occurred when tried to load env variables: %v", err)
	}

	db, err := database.ConnectDatabase()

	if err != nil {
		log.Fatalf("An error occurred when tried to connect to database: %v", err)
	}

	app := SetupApp(db)

	app.Listen(PORT)
}

func SetupApp(db *gorm.DB) *fiber.App {

	app := fiber.New()

	router := app.Group("/api")

	routes.SetupUserRoutes(router, db)

	return app
}
