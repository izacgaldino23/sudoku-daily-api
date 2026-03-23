package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/infrastructure/http/leaderboard"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetLeaderboard(t *testing.T) {
	t.Run("get leaderboard with valid parameters returns entries", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		err = SeedUser("user1@example.com", "player1", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)
		err = SeedUser("user2@example.com", "player2", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)
		err = SeedUser("user3@example.com", "player3", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := GetUserIDByEmail("user1@example.com")
		assert.NoError(t, err)
		user2ID, err := GetUserIDByEmail("user2@example.com")
		assert.NoError(t, err)
		user3ID, err := GetUserIDByEmail("user3@example.com")
		assert.NoError(t, err)

		err = SeedSolve(user1ID, sudokusIDs[0], 30)
		assert.NoError(t, err)
		err = SeedSolve(user2ID, sudokusIDs[0], 60)
		assert.NoError(t, err)
		err = SeedSolve(user3ID, sudokusIDs[0], 45)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response status: %d, body: %s", resp.StatusCode, string(body))

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.Unmarshal(body, &leaderboardResp)
		assert.NoError(t, err)
		assert.Greater(t, len(leaderboardResp.Entries), 0)
		assert.Len(t, leaderboardResp.Entries, 3)
		if len(leaderboardResp.Entries) > 0 {
			assert.Equal(t, 1, leaderboardResp.Entries[0].Rank)
			assert.Equal(t, "player1", leaderboardResp.Entries[0].Username)
			assert.Equal(t, "30", leaderboardResp.Entries[0].Value)
			assert.False(t, leaderboardResp.HasNext)
		}
	})

	t.Run("get leaderboard with no data returns empty list", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			var leaderboardResp leaderboard.LeaderboardResponse
			err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
			assert.NoError(t, err)
			assert.Empty(t, leaderboardResp.Entries)
		} else {
			var errorResponse pkg.Error
			err = json.NewDecoder(resp.Body).Decode(&errorResponse)
			assert.NoError(t, err)
			t.Log(err)
		}
	})

	t.Run("get leaderboard with all-time type returns entries", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		err = SeedUser("user1@example.com", "bestplayer", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := GetUserIDByEmail("user1@example.com")
		assert.NoError(t, err)

		err = SeedSolves(user1ID)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=all-time&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardResp.Entries)
	})

	t.Run("get leaderboard with total type returns entries", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		err = SeedUser("user1@example.com", "activeplayer", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := GetUserIDByEmail("user1@example.com")
		assert.NoError(t, err)

		err = SeedSolves(user1ID)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=total&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardResp.Entries)
	})

	t.Run("get leaderboard with invalid type returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=invalid&size=nine", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with invalid size returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=invalid", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with missing type returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?size=nine", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with missing size returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with limit exceeding max returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=500", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with limit below min returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=0", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with page below min returns bad request", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&page=0", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with pagination returns correct page", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedSudokus()
		assert.NoError(t, err)

		for i := 1; i <= 5; i++ {
			err = SeedUser("user"+string(rune('0'+i))+"@example.com", "player"+string(rune('0'+i)), "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
			assert.NoError(t, err)
			userID, err := GetUserIDByEmail("user" + string(rune('0'+i)) + "@example.com")
			assert.NoError(t, err)
			err = SeedSolve(userID, sudokusIDs[0], i*10)
			assert.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=2&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.Len(t, leaderboardResp.Entries, 2)
		assert.True(t, leaderboardResp.HasNext)
	})

	t.Run("get leaderboard with streak type returns entries", func(t *testing.T) {
		t.Cleanup(TruncateTables)
		app := SetupTestApp()

		err := SeedUser("user1@example.com", "streakplayer", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := GetUserIDByEmail("user1@example.com")
		assert.NoError(t, err)

		ctx := context.Background()
		_, err = database.GetDB().BunConnection.ExecContext(ctx, `
			INSERT INTO user_stats (id, user_id, current_streak, longest_streak, last_solved_date, total_solved)
			VALUES (?, ?, 10, 15, NOW(), 50)
		`, generateUUID(), user1ID)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=streak&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardResp.Entries)
	})
}
