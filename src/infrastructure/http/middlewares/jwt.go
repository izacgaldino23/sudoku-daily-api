package middlewares

import (
	"sudoku-daily-api/src/domain"
	appContext "sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/vo"

	"github.com/gofiber/fiber/v3"
)

const (
	authorizationHeader = "Authorization"
)

func OptionalJWTMiddleware(tokenService domain.TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get(authorizationHeader)

		if len(header) == 0 {
			return c.Next()
		}
		claims, err := tokenService.ParseToken(string(header))
		if err != nil {
			return c.Next()
		}

		userID, ok := claims["user_id"].(vo.UUID)
		if !ok {
			return c.Next()
		}

		// Set userID on context
		reqContext := c.Context()
		newCtx := appContext.SetUserIDOnContext(reqContext, userID)

		c.SetContext(newCtx)

		return c.Next()
	}
}
