package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		reqCtx := c.Context()

		ctx, cancel := context.WithTimeout(reqCtx, timeout)
		defer cancel()

		c.SetContext(ctx)

		return c.Next()
	}
}
