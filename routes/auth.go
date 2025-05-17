package routes

import (
	"todo-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
}

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/users-list", controllers.GetAllUsers)
}
