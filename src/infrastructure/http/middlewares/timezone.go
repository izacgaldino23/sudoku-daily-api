package middlewares

import (
	"github.com/gofiber/fiber/v3"

	"sudoku-daily-api/src/domain/app_context"
)

const (
	TimezoneHeader = "X-Timezone"
)

func NewTimezoneMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		timezone := c.Get(TimezoneHeader)

		if len(timezone) == 0 {
			timezone = "UTC"
		}

		reqCtx := app_context.SetTimezoneOnContext(c.Context(), timezone)

		c.SetContext(reqCtx)

		return c.Next()
	}
}
