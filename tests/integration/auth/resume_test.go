package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestAuthResume(t *testing.T) {
	app := helpers.SetupTestApp()
	t.Cleanup(helpers.TruncateTables)

	setupAuthenticatedUser := func(t *testing.T, withSolves bool) string {
		t.Cleanup(helpers.TruncateTables)

		userData, err := helpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		if withSolves {
			userID, _ := helpers.GetUserIDByEmail(userData.Email)
			err := helpers.SeedSudokus()
			assert.NoError(t, err)

			err = helpers.SeedSolves(userID)
			assert.NoError(t, err)
		}

		return userData.AccessToken
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
	t.Cleanup(helpers.TruncateTables)
	app := helpers.SetupTestApp()

	email := helpers.GenerateUniqueEmail("test")
	tokens, err := helpers.RegisterAndLoginUserWithTokens(app, email, "testuser", "password123")
	assert.NoError(t, err)

	userID, err := helpers.GetUserIDByEmail(email)
	assert.NoError(t, err)

	_ = helpers.SeedSudokus()
	_ = helpers.SeedSolves(userID)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

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
	t.Cleanup(helpers.TruncateTables)
	app := helpers.SetupTestApp()

	tokens1, err := helpers.RegisterAndLoginUserWithTokens(app, "user1@example.com", "user1", "password123")
	assert.NoError(t, err)

	tokens2, err := helpers.RegisterAndLoginUserWithTokens(app, "user2@example.com", "user2", "password123")
	assert.NoError(t, err)

	user1ID, err := helpers.GetUserIDByEmail("user1@example.com")
	assert.NoError(t, err)

	user2ID, err := helpers.GetUserIDByEmail("user2@example.com")
	assert.NoError(t, err)

	_ = helpers.SeedSudokus()
	_ = helpers.SeedSolves(user1ID)
	_ = helpers.SeedSolves(user2ID)

	req1 := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+tokens1.AccessToken)

	resp1, err := app.Test(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	var resumeResp1 auth.ResumeResponse
	err = json.NewDecoder(resp1.Body).Decode(&resumeResp1)
	assert.NoError(t, err)

	req2 := httptest.NewRequest(http.MethodGet, "/api/auth/resume", nil)
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+tokens2.AccessToken)

	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var resumeResp2 auth.ResumeResponse
	err = json.NewDecoder(resp2.Body).Decode(&resumeResp2)
	assert.NoError(t, err)

	assert.Equal(t, resumeResp1.TotalGames, resumeResp2.TotalGames)
	assert.Equal(t, resumeResp1.TodayGames, resumeResp2.TodayGames)
	assert.Equal(t, resumeResp1.BestTimes, resumeResp2.BestTimes)
}
