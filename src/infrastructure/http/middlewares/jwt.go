package middlewares

import (
	"sudoku-daily-api/src/domain"
	appContext "sudoku-daily-api/src/domain/app_context"

	"github.com/gofiber/fiber/v3"
)

const (
	authHeader = "Authorization"
)

type (
	JWTMiddleware interface {
		Execute() func(c fiber.Ctx) error
	}

	jwtMiddlewareImpl struct {
		tokenService domain.TokenService
	}
)

func NewJWTMiddleware(tokenService domain.TokenService) JWTMiddleware {
	return &jwtMiddlewareImpl{
		tokenService: tokenService,
	}
}

func (m *jwtMiddlewareImpl) Execute() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		// Verify if token is present
		header := c.Request().Header.Peek(authHeader)

		// Validate token and get userID
		if len(header) == 0 {
			return fiber.ErrUnauthorized
		}
		userID, err := m.tokenService.ValidateAccessToken(string(header))
		if err != nil {
			return fiber.ErrUnauthorized
		}

		// Set userID on context
		reqContext := c.Context()
		appContext.SetUserOnContext(reqContext, userID)

		// Update context
		c.SetContext(reqContext)

		return c.Next()
	}
}
