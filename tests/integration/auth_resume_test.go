package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/src/infrastructure/http/auth"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthResume(t *testing.T) {
	app := SetupTestApp()

	setupAuthenticatedUser := func(t *testing.T, withSolves bool) string {
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
		_ = json.NewDecoder(loginResp.Body).Decode(&loginResult)

		if withSolves {
			userID, _ := GetUserIDByEmail("test@example.com")
			err := SeedSudokus()
			assert.NoError(t, err)

			err = SeedSolves(userID)
			assert.NoError(t, err)
		}

		return loginResult.AccessToken
	}

	t.Run("valid resume without solves", func(t *testing.T) {
		accessToken := setupAuthenticatedUser(t, false)

		req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var resumeResp auth.ResumeResponse
		err = json.NewDecoder(resp.Body).Decode(&resumeResp)
		assert.NoError(t, err)
		assert.Empty(t, resumeResp.TotalGames)
		assert.Empty(t, resumeResp.TodayGames)
		assert.Empty(t, resumeResp.BestTimes)
	})

	t.Run("valid resume with solves", func(t *testing.T) {
		accessToken := setupAuthenticatedUser(t, true)

		req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var resumeResp auth.ResumeResponse
		err = json.NewDecoder(resp.Body).Decode(&resumeResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resumeResp.TotalGames)
		assert.NotEmpty(t, resumeResp.TodayGames)
		assert.NotEmpty(t, resumeResp.BestTimes)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		_ = setupAuthenticatedUser(t, false)

		req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("invalid access token", func(t *testing.T) {
		_ = setupAuthenticatedUser(t, false)

		req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("malformed authorization header", func(t *testing.T) {
		_ = setupAuthenticatedUser(t, false)

		req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "invalid-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
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

	assert.Equal(t, 3, resumeResp.TotalGames[9])
	assert.Equal(t, 1, resumeResp.TotalGames[4])

	assert.Len(t, resumeResp.TodayGames, 3)

	for _, game := range resumeResp.TodayGames {
		assert.True(t, game.Finished)
		assert.Greater(t, game.Time, 0)
	}

	assert.Len(t, resumeResp.BestTimes, 3)

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
