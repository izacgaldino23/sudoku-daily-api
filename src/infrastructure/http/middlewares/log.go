package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func LogMiddleware(base zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		log := base.With().
			Str("path", c.Path()).
			Str("method", c.Method()).
			Logger()

		ctx := log.WithContext(c.Context())

		c.SetContext(ctx)

		return c.Next()
	}
}
