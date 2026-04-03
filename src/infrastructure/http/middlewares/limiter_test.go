package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"sudoku-daily-api/src/infrastructure/http/middlewares"
)

func TestGlobalRateLimiter(t *testing.T) {
	globalLimit := 10
	app := setupTestAppWithMiddlewares(middlewares.NewGlobalRateLimiterMiddleware(globalLimit))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for i := range globalLimit + 1 {
		response, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		if i < globalLimit {
			assert.Equal(t, http.StatusOK, response.StatusCode)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
		}
	}
}

func TestLocalRateLimiter(t *testing.T) {
	userLimit := 10
	app := setupTestAppWithMiddlewares(middlewares.NewUserRateLimitMiddleware(userLimit))

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for i := range userLimit + 1 {
		response, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)

		if i < userLimit {
			assert.Equal(t, http.StatusOK, response.StatusCode)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
		}
	}
}