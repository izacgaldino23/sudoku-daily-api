package sudoku_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"

	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/tests/integration/testhelpers"
)

func TestSudokuGetMyDailySolves(t *testing.T) {
	t.Run("get solves for today with data returns entries", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		err := testhelpers.SeedSudokus()
		assert.NoError(t, err)

		userData, err := testhelpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, userData.AccessToken)

		userID, err := testhelpers.GetUserIDByEmail(userData.Email)
		assert.NoError(t, err)

		err = testhelpers.SeedSolve(userID, testhelpers.SudokusIDs[0], 60)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer "+userData.AccessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mySolves sudoku.MySolvesResponse
		err = json.NewDecoder(resp.Body).Decode(&mySolves)
		assert.NoError(t, err)
		if assert.Len(t, mySolves.Solves, 1) {
			assert.Equal(t, 60, mySolves.Solves[0].Duration)
		}
	})

	t.Run("get solves for today with no data returns empty list", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		userData, err := testhelpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, userData.AccessToken)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer "+userData.AccessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mySolves sudoku.MySolvesResponse
		err = json.NewDecoder(resp.Body).Decode(&mySolves)
		assert.NoError(t, err)
		assert.Empty(t, mySolves.Solves)
	})

	t.Run("get solves without auth returns 401", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get solves with invalid token returns 401", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("get solves excludes yesterday's solves", func(t *testing.T) {
		t.Cleanup(testhelpers.TruncateTables)
		app := testhelpers.SetupTestApp()

		err := testhelpers.SeedSudokus()
		assert.NoError(t, err)

		userData, err := testhelpers.RegisterAndLoginUser(app, "password123")
		assert.NoError(t, err)

		userID, err := testhelpers.GetUserIDByEmail(userData.Email)
		assert.NoError(t, err)

		solve := testhelpers.SolveSeed{
			ID:        testhelpers.GenerateUUID(),
			UserID:    userID,
			SudokuID:  testhelpers.SudokusIDs[0],
			StartedAt: time.Now().Add(-26 * time.Hour),
			Duration:  90,
			Size:      9,
		}
		ctx := context.Background()
		_, err = database.GetDB().BunConnection.NewInsert().Model(&solve).Exec(ctx)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/sudoku/me", nil)
		req.Header.Set("Authorization", "Bearer "+userData.AccessToken)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var mySolves sudoku.MySolvesResponse
		err = json.NewDecoder(resp.Body).Decode(&mySolves)
		assert.NoError(t, err)
		assert.Empty(t, mySolves.Solves)
	})
}
