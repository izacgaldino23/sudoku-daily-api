package leaderboard_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

type testUserStats struct {
	bun.BaseModel  `bun:"table:user_stats"`
	ID             string    `bun:"id,pk"`
	UserID         string    `bun:"user_id,notnull"`
	CurrentStreak  int       `bun:",notnull"`
	LongestStreak  int       `bun:",notnull"`
	LastSolvedDate time.Time `bun:",notnull"`
	TotalSolved    int       `bun:",notnull"`
}

func resetRequest(t *testing.T, app *fiber.App, date string) *http.Response {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"date": date})
	req := httptest.NewRequest(http.MethodPost, "/api/leaderboard/reset", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	return resp
}

func createUserStats(t *testing.T, userID string, streak int, lastSolvedDate time.Time, totalSolved int) {
	t.Helper()
	db := database.GetDB().BunConnection
	ctx := context.Background()
	_, err := db.NewInsert().Model(&testUserStats{
		ID:             vo.NewUUID().String(),
		UserID:         userID,
		CurrentStreak:  streak,
		LongestStreak:  streak,
		LastSolvedDate: lastSolvedDate,
		TotalSolved:    totalSolved,
	}).Exec(ctx)
	assert.NoError(t, err)
}

func TestResetStrikes(t *testing.T) {
	t.Run("resets streak for stale users and preserves total_solved", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		pass := "strong-password-123"
		userData, err := helpers.RegisterAndLoginUser(app, pass)
		assert.NoError(t, err)

		userID, err := helpers.GetUserIDByEmail(userData.Email)
		assert.NoError(t, err)

		createUserStats(t, userID, 5, time.Now().AddDate(0, 0, -3), 3)

		resp := resetRequest(t, app, time.Now().Format(time.RFC3339))
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		solves, streak, err := helpers.GetUserStats(userID)
		assert.NoError(t, err)
		assert.Equal(t, 0, streak, "current_streak should be reset to 0")
		assert.Equal(t, 3, solves, "total_solved should be preserved")
	})

	t.Run("preserves streak for active users with recent last_solved_date", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		pass := "strong-password-123"
		userData, err := helpers.RegisterAndLoginUser(app, pass)
		assert.NoError(t, err)

		userID, err := helpers.GetUserIDByEmail(userData.Email)
		assert.NoError(t, err)

		createUserStats(t, userID, 5, time.Now().AddDate(0, 0, -1), 10)

		resp := resetRequest(t, app, time.Now().Format(time.RFC3339))
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, streak, err := helpers.GetUserStats(userID)
		assert.NoError(t, err)
		assert.Equal(t, 5, streak, "active user's streak should be preserved")
	})

	t.Run("mixed users: only stale users get reset", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		user1, err := helpers.RegisterAndLoginUser(app, "strong-password-123")
		assert.NoError(t, err)
		user2, err := helpers.RegisterAndLoginUser(app, "strong-password-456")
		assert.NoError(t, err)

		uid1, err := helpers.GetUserIDByEmail(user1.Email)
		assert.NoError(t, err)
		uid2, err := helpers.GetUserIDByEmail(user2.Email)
		assert.NoError(t, err)

		createUserStats(t, uid1, 5, time.Now().AddDate(0, 0, -3), 1)
		createUserStats(t, uid2, 3, time.Now().AddDate(0, 0, -1), 5)

		resp := resetRequest(t, app, time.Now().Format(time.RFC3339))
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		_, streak1, err := helpers.GetUserStats(uid1)
		assert.NoError(t, err)
		assert.Equal(t, 0, streak1, "stale user's streak should be reset")

		_, streak2, err := helpers.GetUserStats(uid2)
		assert.NoError(t, err)
		assert.Equal(t, 3, streak2, "active user's streak should be preserved")
	})

	t.Run("returns 204 when no users exist", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		resp := resetRequest(t, app, time.Now().Format(time.RFC3339))
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("returns 204 when date field is missing (defaults to now)", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/leaderboard/reset", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("returns 204 when body is empty", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/leaderboard/reset", http.NoBody)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("returns 400 when date format is invalid", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/leaderboard/reset", strings.NewReader(`{"date": "not-a-date"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("returns 400 when body is not json", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/leaderboard/reset", strings.NewReader("not-json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("idempotent: calling reset twice has same effect", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		pass := "strong-password-123"
		userData, err := helpers.RegisterAndLoginUser(app, pass)
		assert.NoError(t, err)

		userID, err := helpers.GetUserIDByEmail(userData.Email)
		assert.NoError(t, err)

		createUserStats(t, userID, 5, time.Now().AddDate(0, 0, -3), 2)

		resp1 := resetRequest(t, app, time.Now().Format(time.RFC3339))
		assert.Equal(t, http.StatusNoContent, resp1.StatusCode)

		resp2 := resetRequest(t, app, time.Now().Format(time.RFC3339))
		assert.Equal(t, http.StatusNoContent, resp2.StatusCode)

		_, streak, err := helpers.GetUserStats(userID)
		assert.NoError(t, err)
		assert.Equal(t, 0, streak)
	})
}
