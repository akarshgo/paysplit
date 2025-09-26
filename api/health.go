package api

import "github.com/gofiber/fiber/v2"

//HanldeHealth returns a simple ok response
//useful for uptime checks, monitoring and debugging

func HandleHealth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
