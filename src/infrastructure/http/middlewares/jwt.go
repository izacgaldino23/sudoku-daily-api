package middlewares

import (
	"sudoku-daily-api/src/domain"
	appContext "sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/logging"

	"github.com/gofiber/fiber/v3"
)

const (
	authorizationHeader = "Authorization"
)

func OptionalJWTMiddleware(tokenService domain.TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		reqContext := c.Context()
		logger := logging.Log(reqContext)

		if appContext.GetUserIDFromContext(c.Context()) != "" {
			logger.Info().Msg("user already authenticated")
			return c.Next()
		}

		authHeader := string(c.Request().Header.Peek(authorizationHeader))

		if authHeader == "" {
			logger.Info().Msg("no token provided")
			return c.Next()
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			logger.Info().Msg("token with Bearer")
			tokenString = authHeader[7:]
		}

		claims, err := tokenService.ParseToken(tokenString)
		if err != nil {
			logger.Info().Err(err).Msg("invalid token")
			return c.Next()
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.Info().Err(err).Msg("user_id not found in claims")
			return c.Next()
		}

		// Set userID on context
		newCtx := appContext.SetUserIDOnContext(reqContext, vo.UUID(userID))

		c.SetContext(newCtx)

		return c.Next()
	}
}
