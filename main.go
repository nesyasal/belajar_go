package main

import (
    "todo-app/database"
    "todo-app/routes"

    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
    database.Connect()
	app.Use(logger.New(logger.Config{
        Format: "${status} - ${method} ${path}\n",
    }))
    routes.AuthRoutes(app)
    routes.SetupRoutes(app)

    app.Listen(":3000")
}
