package middlewares

import (
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
)

func NewRequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")

		if len(requestID) == 0 {
			requestID = vo.NewUUID().String()
		}

		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}
