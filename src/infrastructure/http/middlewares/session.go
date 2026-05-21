package middlewares

import (
	"sudoku-daily-api/pkg"
	appContext "sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
)

const (
	XSessionIdHeader = "X-Session-Id"
)

func SessionMiddleware() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		var (
			header    = c.Get(XSessionIdHeader)
			sessionID = newSessionID()
		)

		if len(header) > 0 {
			sessionID = vo.UUID(header)

			if !sessionID.IsValid() {
				return pkg.JsonError(c, pkg.ErrInvalidToken)
			}
		}

		reqContext := c.Context()
		newCtx := appContext.SetSessionIDOnContext(reqContext, sessionID)

		c.SetContext(newCtx)

		return c.Next()
	}
}

func newSessionID() vo.UUID {
	return vo.NewUUID()
}
