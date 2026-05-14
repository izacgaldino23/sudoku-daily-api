package middlewares

import (
	"strings"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/logging"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/api/idtoken"
)

const cronSecretHeader = "X-Cron-Secret"

func AuthOIDCMiddleware(isEnabled bool, audience, cronSecret string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		reqCtx := c.Context()

		if cronSecret != "" {
			if c.Get(cronSecretHeader) == cronSecret {
				return c.Next()
			}
			authHeader := c.Get(authorizationHeader)
			if strings.TrimPrefix(authHeader, "Bearer ") == cronSecret {
				return c.Next()
			}
		}

		if !isEnabled {
			return pkg.JsonError(c, pkg.ErrInvalidToken)
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
