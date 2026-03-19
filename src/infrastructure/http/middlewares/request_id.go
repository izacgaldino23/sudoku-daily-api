package middlewares

import (
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
)

const (
	XRequestIDHeader = "X-Request-ID"
)

func NewRequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get(XRequestIDHeader)

		if len(requestID) == 0 {
			requestID = vo.NewUUID().String()
		}

		reqCtx := app_context.SetRequestIDOnContext(c.Context(), vo.UUID(requestID))

		c.SetContext(reqCtx)

		return c.Next()
	}
}
