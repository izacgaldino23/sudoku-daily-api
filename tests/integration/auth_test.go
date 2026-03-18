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

func TestAuthResume(t *testing.T) {
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
		setupData  func()
		wantStatus int
		checkResp  func(*testing.T, *http.Response)
	}{
		{
			name:       "valid resume without solves",
			headers:    map[string]string{"Authorization": "Bearer " + loginResult.AccessToken},
			setupData:  func() {},
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, resp *http.Response) {
				var resumeResp auth.ResumeResponse
				err := json.NewDecoder(resp.Body).Decode(&resumeResp)
				assert.NoError(t, err)
				assert.Empty(t, resumeResp.TotalGames)
				assert.Empty(t, resumeResp.TodayGames)
				assert.Empty(t, resumeResp.BestTimes)
			},
		},
		{
			name:       "valid resume with solves",
			headers:    map[string]string{"Authorization": "Bearer " + loginResult.AccessToken},
			setupData:  func() { _ = SeedSudokus(); userID, _ := GetUserIDByEmail("test@example.com"); _ = SeedSolves(userID) },
			wantStatus: http.StatusOK,
			checkResp: func(t *testing.T, resp *http.Response) {
				var resumeResp auth.ResumeResponse
				err := json.NewDecoder(resp.Body).Decode(&resumeResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resumeResp.TotalGames)
				assert.NotEmpty(t, resumeResp.TodayGames)
				assert.NotEmpty(t, resumeResp.BestTimes)
			},
		},
		{
			name:       "missing authorization header",
			headers:    map[string]string{},
			setupData:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid access token",
			headers:    map[string]string{"Authorization": "Bearer invalid-token"},
			setupData:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed authorization header",
			headers:    map[string]string{"Authorization": "invalid-token"},
			setupData:  func() {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
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

			tt.setupData()

			req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
			req.Header.Set("Content-Type", "application/json")

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			resp, err := app.Test(req, fiber.TestConfig{
				Timeout: 0,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.checkResp != nil && resp.StatusCode == http.StatusOK {
				tt.checkResp(t, resp)
			}
		})
	}
}

func TestAuthResume_VerifyDataAccuracy(t *testing.T) {
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

	userID, err := GetUserIDByEmail("test@example.com")
	assert.NoError(t, err)

	_ = SeedSudokus()
	_ = SeedSolves(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResult.AccessToken)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var resumeResp auth.ResumeResponse
	err = json.NewDecoder(resp.Body).Decode(&resumeResp)
	assert.NoError(t, err)

	assert.Equal(t, 2, resumeResp.TotalGames[9])
	assert.Equal(t, 1, resumeResp.TotalGames[4])

	assert.Len(t, resumeResp.TodayGames, 3)

	for _, game := range resumeResp.TodayGames {
		assert.True(t, game.Finished)
		assert.Greater(t, game.Time, 0)
	}

	assert.Len(t, resumeResp.BestTimes, 2)

	for _, game := range resumeResp.BestTimes {
		assert.True(t, game.Finished)
		assert.Greater(t, game.Time, 0)
	}

	if len(resumeResp.BestTimes) > 0 {
		size9Best := 0
		for _, game := range resumeResp.BestTimes {
			if game.Size == 9 {
				size9Best = game.Time
				break
			}
		}
		assert.Equal(t, 60, size9Best, "best time for size 9 should be 60 seconds (fastest solve)")
	}
}

func TestAuthResume_MultipleUsers(t *testing.T) {
	TruncateTables(t)

	app := SetupTestApp()

	registerBody1, _ := json.Marshal(map[string]string{
		"email":    "user1@example.com",
		"username": "user1",
		"password": "password123",
	})
	registerReq1 := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody1))
	registerReq1.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(registerReq1)

	registerBody2, _ := json.Marshal(map[string]string{
		"email":    "user2@example.com",
		"username": "user2",
		"password": "password123",
	})
	registerReq2 := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(registerBody2))
	registerReq2.Header.Set("Content-Type", "application/json")
	_, _ = app.Test(registerReq2)

	loginBody1, _ := json.Marshal(map[string]string{
		"email":    "user1@example.com",
		"password": "password123",
	})
	loginReq1 := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody1))
	loginReq1.Header.Set("Content-Type", "application/json")
	loginResp1, _ := app.Test(loginReq1)

	loginBody2, _ := json.Marshal(map[string]string{
		"email":    "user2@example.com",
		"password": "password123",
	})
	loginReq2 := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody2))
	loginReq2.Header.Set("Content-Type", "application/json")
	loginResp2, _ := app.Test(loginReq2)

	var loginResult1, loginResult2 auth.LoginResponse
	_ = json.NewDecoder(loginResp1.Body).Decode(&loginResult1)
	_ = json.NewDecoder(loginResp2.Body).Decode(&loginResult2)

	_ = SeedSudokus()

	user1ID, _ := GetUserIDByEmail("user1@example.com")
	user2ID, _ := GetUserIDByEmail("user2@example.com")

	_ = SeedSolves(user1ID)
	_ = SeedSolves(user2ID)

	req1 := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+loginResult1.AccessToken)

	req2 := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+loginResult2.AccessToken)

	resp1, _ := app.Test(req1)
	resp2, _ := app.Test(req2)

	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var resumeResp1, resumeResp2 auth.ResumeResponse
	_ = json.NewDecoder(resp1.Body).Decode(&resumeResp1)
	_ = json.NewDecoder(resp2.Body).Decode(&resumeResp2)

	assert.Equal(t, resumeResp1.TotalGames, resumeResp2.TotalGames)

	assert.Len(t, resumeResp1.TodayGames, 3)
	assert.Len(t, resumeResp2.TodayGames, 3)
}
