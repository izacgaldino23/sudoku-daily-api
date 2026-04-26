package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
)

func TestAuthLogout(t *testing.T) {
	setupTokens := func() (string, string) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		email := helpers.GenerateUniqueEmail("test")
		username := helpers.GenerateUniqueUsername("testuser")
		tokens, err := helpers.RegisterAndLoginUserWithTokens(app, email, username, "password123")
		assert.NoError(t, err)

		return tokens.AccessToken, tokens.RefreshToken
	}

	t.Run("valid logout", func(t *testing.T) {
		accessToken, refreshToken := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)
		helpers.SetRefreshCookie(req, refreshToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookies := resp.Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "refresh_token", cookies[0].Name)
		assert.Equal(t, "", cookies[0].Value)
		assert.True(t, cookies[0].MaxAge < 0)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		_, refreshToken := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		helpers.SetRefreshCookie(req, refreshToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid access token", func(t *testing.T) {
		_, refreshToken := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid-token")
		helpers.SetRefreshCookie(req, refreshToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("missing refresh token cookie", func(t *testing.T) {
		accessToken, _ := setupTokens()
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var result pkg.Error
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.Message)
	})
}