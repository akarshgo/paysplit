package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

const (
	listeningPort = ":3000"
)

func main() {
	fmt.Println("Paysplit initialised")

	app := fiber.New()

	app.Get("/", handleFoo)

	log.Fatal(app.Listen(listeningPort))
}

func handleFoo(c *fiber.Ctx) error {
	return c.SendString("Hello from foo")
}
