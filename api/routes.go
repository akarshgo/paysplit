package api

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, userHandlers *UserHandlers) {
	// v1 prefix
	v1 := app.Group("/v1")

	// Users
	v1.Post("/users", userHandlers.HandleCreateUser)
	v1.Get("/users", userHandlers.HandleGetUsers)
	v1.Get("/users/:id", userHandlers.HandleGetUser)
	v1.Patch("/users/:id", userHandlers.HandleUpdateUser)
	v1.Delete("/users/:id", userHandlers.HandleDeleteUser)

	// Later: add Groups, Expenses, etc.
}
