package middlewares

import (
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

		c.Set(XRequestIDHeader, requestID)
		return c.Next()
	}
}