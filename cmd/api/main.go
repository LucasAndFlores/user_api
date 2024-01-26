package main

import (
	"fmt"
	"os"

	"github.com/LucasAndFlores/user_api/config"
	"github.com/gofiber/fiber/v2"
)


var PORT = fmt.Sprintf(":%v", os.Getenv("PORT"))

func main() {
    err := config.LoadEnvVariables()
    
    if err != nil {
        panic(err)
    }

    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("ok")
    })

    
    app.Listen(PORT)
}
