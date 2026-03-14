package middlewares

import (
	"sudoku-daily-api/src/domain/app_context"

	"github.com/gofiber/fiber/v3"
)

func NewResponseHeadersMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()

		reqCtx := c.Context()
		sessionID := app_context.GetSessionIDFromContext(reqCtx)
		c.Set("X-Session-ID", sessionID.String())

		requestID := app_context.GetRequestIDFromContext(reqCtx)
		c.Set("X-Request-ID", requestID.String())

		return err
	}
}
