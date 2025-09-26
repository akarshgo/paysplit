package main

import (
	"database/sql"
	"log"

	"github.com/akarshgo/paysplit/api"
	"github.com/akarshgo/paysplit/db"
	"github.com/akarshgo/paysplit/logger"
	rediscli "github.com/akarshgo/paysplit/redis"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	logger.Init()
	defer logger.Sync()

	dsn := "postgres://paysplit:paysplit@localhost:5432/paysplit?sslmode=disable"
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	userStore := db.NewPostgresUserStore(sqlDB)
	userHandlers := api.NewUserHandlers(userStore)
	groupStore := db.NewPostgresGroupStore(sqlDB)
	groupHandlers := api.NewGroupHanlders(groupStore)
	expenseStore := db.NewPostgresExpenseStore(sqlDB)
	expenseHandlers := api.NewExpenseHandlers(expenseStore)
	linkHanlders := api.NewLinksHandlers("paysplit")

	// Initialize Redis before starting the HTTP server
	rediscli.Init()
	defer func() {
		_ = rediscli.Rdb.Close()
	}()

	app := fiber.New()
	api.SetupRoutes(app, userHandlers, groupHandlers, expenseHandlers, linkHanlders)

	log.Println("API on :8080")
	app.Listen(":8080")
}
