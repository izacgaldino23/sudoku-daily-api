package middlewares

import (
	"sudoku-daily-api/src/domain"
	appContext "sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
)

const (
	sessionHeader = "session"
)

func SessionMiddleware(tokenService domain.TokenService) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		header := c.Get(sessionHeader)

		if len(header) == 0 {
			return c.Next()
		}
		claims, err := tokenService.ParseToken(string(header))
		if err != nil {
			return c.Next()
		}

		sessionID, ok := claims["session_id"].(vo.UUID)
		if !ok {
			return c.Next()
		}

		reqContext := c.Context()
		newCtx := appContext.SetSessionIDOnContext(reqContext, sessionID)

		c.SetContext(newCtx)

		return c.Next()
	}
}
