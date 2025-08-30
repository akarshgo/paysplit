package api

import (
	"github.com/akarshgo/paysplit/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := uuid.New().String()
		c.Set("X-Request-ID", reqID)
		c.Locals("request_id", reqID)

		// log each request
		logger.Log.Info("incoming request",
			zap.String("id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		return c.Next()
	}
}
