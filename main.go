package main

import (
	"database/sql"
	"log"

	"github.com/akarshgo/paysplit/api"
	"github.com/akarshgo/paysplit/db"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://paysplit:paysplit@localhost:5432/paysplit?sslmode=disable"
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	userStore := db.NewPostgresUserStore(sqlDB)
	userHandlers := api.NewUserHandlers(userStore)
	groupStore := db.NewPostGresGroupStore(sqlDB)
	groupHandlers := api.NewGroupHanlders(groupStore)
	expenseStore := db.NewPostGresExpenseStore(sqlDB)
	expenseHandlers := api.NewExpenseHandlers(expenseStore)

	app := fiber.New()
	api.SetupRoutes(app, userHandlers, groupHandlers, expenseHandlers)

	log.Println("API on :8080")
	app.Listen(":8080")
}
