package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/middlewares"
)

func TestLogMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		body          string
		statusCode    int
		responseBody  interface{}
		errorToReturn error
	}{
		{
			name:       "Should log POST request with body",
			method:     http.MethodPost,
			body:       `{"name":"test"}`,
			statusCode: http.StatusOK,
		},
		{
			name:          "Should log error from handler returning error",
			method:        http.MethodPost,
			body:          "",
			statusCode:    http.StatusInternalServerError,
			errorToReturn: pkg.ErrSudokuNotFound,
		},
		{
			name:         "Should log error when status >= 300",
			method:       http.MethodGet,
			statusCode:   http.StatusNotFound,
			responseBody: pkg.NewError("not found"),
		},
		{
			name:         "Should handle invalid response body on error",
			method:       http.MethodGet,
			statusCode:   http.StatusBadRequest,
			responseBody: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := zerolog.New(buf)

			app := setupTestAppWithMiddlewares(middlewares.NewRequestIDMiddleware(), middlewares.LogMiddleware(logger))

			handler := func(c fiber.Ctx) error {
				if tt.errorToReturn != nil {
					return tt.errorToReturn
				}

				if tt.responseBody != nil {
					return c.Status(tt.statusCode).JSON(tt.responseBody)
				}
				return c.SendStatus(tt.statusCode)
			}

			switch tt.method {
			case http.MethodGet:
				app.Get("/test", handler)
			case http.MethodPost:
				app.Post("/test", handler)
			}

			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, "/test", bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			}

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, resp.StatusCode)

			logOutput := buf.String()
			assert.NotEmpty(t, logOutput, "expected some logging to occur")
			assert.Contains(t, logOutput, "/test")
			assert.Contains(t, logOutput, tt.method)

			if tt.method != http.MethodGet && tt.body != "" {
				assert.Contains(t, logOutput, tt.body)
			}

			if tt.errorToReturn != nil {
				assert.Contains(t, logOutput, "error")
			}

			if tt.responseBody != nil && tt.statusCode >= 300 {
				if errMsg, ok := tt.responseBody.(pkg.Error); ok {
					assert.Contains(t, logOutput, errMsg.Message)
				}
			}
		})
	}
}

func TestLogMiddleware_GETNoError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	app := setupTestAppWithMiddlewares(middlewares.NewRequestIDMiddleware(), middlewares.LogMiddleware(logger))
	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	resp, err := app.Test(req, fiber.TestConfig{
		Timeout: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	logOutput := buf.String()
	assert.NotContains(t, "error", logOutput)
}

func TestLogMiddleware_JsonErrorResponse(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	app := setupTestAppWithMiddlewares(middlewares.NewRequestIDMiddleware(), middlewares.LogMiddleware(logger))
	app.Get("/error", func(c fiber.Ctx) error {
		return c.Status(http.StatusBadRequest).JSON(pkg.Error{Message: "validation failed"})
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)

	resp, err := app.Test(req, fiber.TestConfig{
		Timeout: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "validation failed")
}

func TestLogMiddleware_EmptyResponseBody(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	app := setupTestAppWithMiddlewares(middlewares.NewRequestIDMiddleware(), middlewares.LogMiddleware(logger))
	app.Get("/empty", func(c fiber.Ctx) error {
		return c.Status(http.StatusNoContent).Send(nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/empty", nil)

	_, err := app.Test(req, fiber.TestConfig{
		Timeout: 0,
	})
	assert.NoError(t, err)

	logOutput := buf.String()
	assert.NotContains(t, "error", logOutput)
}

func TestLogMiddleware_RequestIDFromContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	app := setupTestAppWithMiddlewares(middlewares.NewRequestIDMiddleware(), middlewares.LogMiddleware(logger))
	app.Post("/test", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"test":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", "custom-request-id")

	_, err := app.Test(req, fiber.TestConfig{
		Timeout: 0,
	})
	assert.NoError(t, err)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "custom-request-id")
}

func TestLogMiddleware_GlobalLoggerNotUsed(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)

	app := setupTestAppWithMiddlewares(middlewares.NewRequestIDMiddleware(), middlewares.LogMiddleware(logger))
	app.Post("/test", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"test":1}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, fiber.TestConfig{
		Timeout: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("Log output: %s", buf.String())

	// buf (passed logger) should be used
	assert.Contains(t, buf.String(), "/test")
}

func setupTestAppWithMiddlewares(middlewaresList ...fiber.Handler) *fiber.App {
	app := fiber.New()
	for _, middleware := range middlewaresList {
		app.Use(middleware)
	}
	// app.Use(middlewares.NewRequestIDMiddleware())
	// middleware := middlewares.LogMiddleware(logger)
	// app.Use(middleware)

	return app
}
