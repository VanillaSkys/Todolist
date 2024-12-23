package middleware

import (
	"github.com/VanillaSkys/todo_fiber/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func SetRequestId() fiber.Handler {
	return func(c fiber.Ctx) error {
		reqId := c.Get("X-Request-ID", "")

		if reqId == "" {
			reqId = uuid.New().String()
		}

		loggerWithReqID := logger.AddRequestIDToLogger(reqId)
		c.Locals("X-Request-ID", reqId)
		loggerWithReqID.Info("Incoming request", zap.String("method", c.Method()), zap.String("path", c.Path()))

		c.Locals("logger", loggerWithReqID)

		err := c.Next()

		loggerWithReqID.Info("Request completed", zap.Int("status", c.Response().StatusCode()))

		return err
	}
}
