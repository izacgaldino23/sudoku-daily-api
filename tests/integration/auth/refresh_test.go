package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthRefresh(t *testing.T) {
	setupTokens := func() (string, string) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		email := helpers.GenerateUniqueEmail("test")
		username := helpers.GenerateUniqueUsername("testuser")
		tokens, err := helpers.RegisterAndLoginUserWithTokens(app, email, username, "password123")
		assert.NoError(t, err)

		return tokens.AccessToken, tokens.RefreshToken
	}

	t.Run("valid refresh", func(t *testing.T) {
		accessToken, refreshToken := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		helpers.SetRefreshCookie(req, refreshToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var refreshResp auth.RefreshTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&refreshResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshResp.AccessToken)

		cookies := resp.Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "refresh_token", cookies[0].Name)
		assert.NotEmpty(t, cookies[0].Value)
	})

	t.Run("missing refresh token cookie", func(t *testing.T) {
		accessToken, _ := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var result pkg.Error
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("invalid refresh token cookie", func(t *testing.T) {
		accessToken, _ := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		helpers.SetRefreshCookie(req, "invalid-refresh-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var result pkg.Error
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Message)
	})
}