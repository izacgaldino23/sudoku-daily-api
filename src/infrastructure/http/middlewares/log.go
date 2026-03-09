package middlewares

import (
	"context"
	"sudoku-daily-api/src/infrastructure/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func LogMiddleware(base zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {

		log := base.With().
			Str("path", c.Path()).
			Str("method", c.Method()).
			Logger()

		ctx := context.WithValue(c.Context(), logging.LoggerKey, &log)

		c.SetContext(ctx)

		return c.Next()
	}
}
