package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthRegister(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "valid registration",
			body:       map[string]string{"email": "test@example.com", "username": "testuser", "password": "password123"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing email",
			body:       map[string]string{"username": "testuser", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email format",
			body:       map[string]string{"email": "invalid-email", "username": "testuser", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing username",
			body:       map[string]string{"email": "test@example.com", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "username too short",
			body:       map[string]string{"email": "test@example.com", "username": "ab", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       map[string]string{"email": "test@example.com", "username": "testuser"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "password too short",
			body:       map[string]string{"email": "test@example.com", "username": "testuser", "password": "123"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TruncateTables(t)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			// expected body response here is empty
			if tt.wantStatus != resp.StatusCode {
				result := map[string]interface{}{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Empty(t, result)
			}
		})
	}
}

func TestAuthLogin(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	registerBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"username": "testuser",
		"password": "password123",
	})
	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(registerReq)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "valid login",
			body:       map[string]string{"email": "test@example.com", "password": "password123"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong password",
			body:       map[string]string{"email": "test@example.com", "password": "wrongpassword"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "user not found",
			body:       map[string]string{"email": "nonexistent@example.com", "password": "password123"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing email",
			body:       map[string]string{"password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       map[string]string{"email": "test@example.com"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TruncateTables(t)

			// Register
			registerBody, _ := json.Marshal(map[string]string{
				"email":    "test@example.com",
				"username": "testuser",
				"password": "password123",
			})
			registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody))
			registerReq.Header.Set("Content-Type", "application/json")
			_, _ = app.Test(registerReq)

			// Login
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var loginResp auth.LoginResponse
				err := json.NewDecoder(resp.Body).Decode(&loginResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, loginResp.AccessToken)
				assert.NotEmpty(t, loginResp.RefreshToken)
				assert.Equal(t, "testuser", loginResp.UserName)
				assert.Equal(t, "test@example.com", loginResp.Email)
			} else {
				result := pkg.Error{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)

				if tt.wantStatus != http.StatusOK {
					assert.NotEmpty(t, result.Message)
				} else {
					assert.Empty(t, result.Message)
				}
			}
		})
	}
}

func TestAuthRefresh(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	// register
	registerBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"username": "testuser",
		"password": "password123",
	})
	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(registerReq)

	// login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := app.Test(loginReq)

	var loginResult auth.LoginResponse
	err := json.NewDecoder(loginResp.Body).Decode(&loginResult)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		headers    map[string]string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "valid refresh",
			headers:    map[string]string{"Authorization": loginResult.AccessToken},
			body:       map[string]string{"refresh_token": loginResult.RefreshToken},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing refresh token",
			headers:    map[string]string{"Authorization": loginResult.AccessToken},
			body:       map[string]string{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid refresh token",
			headers:    map[string]string{"Authorization": loginResult.AccessToken},
			body:       map[string]string{"refresh_token": "invalid-refresh-token"},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Refresh
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var refreshResp auth.RefreshTokenResponse
				err := json.NewDecoder(resp.Body).Decode(&refreshResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, refreshResp.AccessToken)
			} else {
				result := pkg.Error{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				assert.NoError(t, err)

				if tt.wantStatus != http.StatusOK {
					assert.NotEmpty(t, result.Message)
				} else {
					assert.Empty(t, result.Message)
				}
			}
		})
	}
}

func TestAuthLogout(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	registerBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"username": "testuser",
		"password": "password123",
	})
	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(registerReq)

	loginBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := app.Test(loginReq)

	var loginResult auth.LoginResponse
	err := json.NewDecoder(loginResp.Body).Decode(&loginResult)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		headers    map[string]string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "valid logout",
			headers:    map[string]string{"Authorization": "Bearer " + loginResult.AccessToken},
			body:       map[string]string{"refresh_token": loginResult.RefreshToken},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing authorization header",
			headers:    map[string]string{},
			body:       map[string]string{"refresh_token": loginResult.RefreshToken},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid access token",
			headers:    map[string]string{"Authorization": "invalid-token"},
			body:       map[string]string{"refresh_token": loginResult.RefreshToken},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing refresh token",
			headers:    map[string]string{"Authorization": loginResult.AccessToken},
			body:       map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TruncateTables(t)

			registerBody, _ := json.Marshal(map[string]string{
				"email":    "test@example.com",
				"username": "testuser",
				"password": "password123",
			})
			registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody))
			registerReq.Header.Set("Content-Type", "application/json")
			_, _ = app.Test(registerReq)

			loginBody, _ := json.Marshal(map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			})
			loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
			loginReq.Header.Set("Content-Type", "application/json")
			loginResp, _ := app.Test(loginReq)

			var loginResult auth.LoginResponse
			err := json.NewDecoder(loginResp.Body).Decode(&loginResult)
			assert.NoError(t, err)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
