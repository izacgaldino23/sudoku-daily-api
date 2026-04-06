package middlewares

import (
	"strings"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	appContext "sudoku-daily-api/src/domain/app_context"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func RequireJWTMiddleware(tokenService domain.TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		reqContext := c.Context()
		logger := log.Ctx(reqContext)

		if appContext.GetUserIDFromContext(c.Context()) != "" {
			return c.Next()
		}

		// Verify if token is present
		header := string(c.Request().Header.Peek(authorizationHeader))

		// Validate token and get userID
		if len(header) == 0 {
			logger.Warn().Str("header", header).Msg("empty token provided")
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		}

		// verify if Bearer is passed and remove
		if !strings.HasPrefix(header, "Bearer") {
			logger.Warn().Str("header", header).Msg("token without Bearer")
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		} else if parts := strings.Split(header, " "); len(parts) != 2 {
			logger.Warn().Str("header", header).Msg("token with invalid format")
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		} else {
			header = parts[1]
		}

		userID, err := tokenService.ValidateAccessToken(string(header))
		if err != nil {
			logger.Error().Err(err).Msg("error validating token")
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		}

		if !userID.IsValid() {
			logger.Error().Err(err).Msg("invalid userID uuid")
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		}

		// Set userID on context
		newCtx := appContext.SetUserIDOnContext(reqContext, userID)

		c.SetContext(newCtx)

		return c.Next()
	}
}

func AuthMinimumMiddleware(tokenService domain.TokenService) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		sessionID := appContext.GetSessionIDFromContext(c.Context())
		userID := appContext.GetUserIDFromContext(c.Context())

		if sessionID == "" && userID == "" {
			return pkg.JsonError(c, pkg.ErrInvalidToken)
		}

		return c.Next()
	}
}
