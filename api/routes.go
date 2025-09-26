package api

import (
	"context"

	rediscli "github.com/akarshgo/paysplit/redis"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandlers *UserHandlers, groupHandlers *GroupHandlers, expenseHandlers *ExpenseHandlers, linksHandlers *LinksHandlers) {
	// v1 prefix
	v1 := app.Group("/v1")

	// Users
	v1.Post("/users", userHandlers.HandleCreateUser)
	v1.Get("/users", userHandlers.HandleGetUsers)
	v1.Get("/users/:id", userHandlers.HandleGetUser)
	v1.Patch("/users/:id", userHandlers.HandleUpdateUser)
	v1.Delete("/users/:id", userHandlers.HandleDeleteUser)

	// Groups
	v1.Post("/groups", groupHandlers.HandleCreateGroup)
	v1.Post("/groups/:id/members", groupHandlers.HandleAddMember)
	v1.Get("/groups", groupHandlers.HandleListGroups)

	//Expenses
	v1.Post("/groups/:id/expenses", expenseHandlers.HandleCreateExpense)
	v1.Get("/groups/:id/expenses", expenseHandlers.HandleListExpenses)

	v1.Get("/groups/:id/balances", expenseHandlers.HandleGroupBalances)
	v1.Get("/groups/:id/simplify", expenseHandlers.HandleSimplifyDebts)

	//UPI Links
	v1.Post("/links/settle", linksHandlers.HandleBuildSettleLink)

	//Helath Check
	v1.Get("/health", HandleHealth)

	//Redis
	v1.Get("/ping-redis", func(c *fiber.Ctx) error {
		ctx := context.Background()
		if err := rediscli.Rdb.Set(ctx, "paysplit:ping", "pong", 0).Err(); err != nil {
			return c.Status(500).JSON(fiber.Map{"redis": "error", "err": err.Error()})
		}
		val, err := rediscli.Rdb.Get(ctx, "paysplit:ping").Result()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"redis": "error", "err": err.Error()})
		}
		return c.JSON(fiber.Map{"redis": "ok", "val": val})
	})

}
