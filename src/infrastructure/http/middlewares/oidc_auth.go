package middlewares

import (
	"strings"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/logging"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/api/idtoken"
)

func AuthOIDCMiddleware(isEnabled bool, audience string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		reqCtx := c.Context()

		if !isEnabled {
			return c.Next()
		}

		authHeader := c.Get(authorizationHeader)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := idtoken.Validate(reqCtx, token, audience)
		if err != nil {
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		}

		logging.Log(reqCtx).Info().Str("sub", payload.Subject).Msg("authenticated")

		return c.Next()
	}
}
