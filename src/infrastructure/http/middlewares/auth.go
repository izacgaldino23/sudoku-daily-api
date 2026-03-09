package middlewares

import (
	"net/http"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	appContext "sudoku-daily-api/src/domain/app_context"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(tokenService domain.TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if appContext.GetUserIDFromContext(c.Context()) != "" {
			return c.Next()
		}

		// Verify if token is present
		header := c.Get(authHeader)

		// Validate token and get userID
		if len(header) == 0 {
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		}
		userID, err := tokenService.ValidateAccessToken(string(header))
		if err != nil {
			return pkg.JsonErrorWithStatus(c, err, http.StatusUnauthorized)
		}

		// Set userID on context
		reqContext := c.Context()
		newCtx := appContext.SetUserOnContext(reqContext, userID)

		c.SetContext(newCtx)

		return c.Next()
	}
}
