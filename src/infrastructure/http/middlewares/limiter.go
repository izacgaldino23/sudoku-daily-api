package middlewares

import (
	"time"

	"sudoku-daily-api/pkg"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func NewGlobalRateLimiterMiddleware(limit int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max: limit,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return "global"
		},
		LimitReached: func(c fiber.Ctx) error {
			return pkg.JsonError(c, pkg.ErrTooManyRequests)
		},
	})
}

func NewUserRateLimitMiddleware(limit int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max: limit,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return pkg.JsonError(c, pkg.ErrTooManyRequests)
		},
	})
}