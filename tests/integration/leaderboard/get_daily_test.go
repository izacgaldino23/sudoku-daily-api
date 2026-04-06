package leaderboard_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/src/infrastructure/http/leaderboard"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
)

func TestGetDailyLeaderboard(t *testing.T) {
	t.Run("get leaderboard with valid parameters returns entries", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		email1 := helpers.GenerateUniqueEmail("user1")
		email2 := helpers.GenerateUniqueEmail("user2")
		email3 := helpers.GenerateUniqueEmail("user3")

		err = helpers.SeedUser(email1, "player1", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)
		err = helpers.SeedUser(email2, "player2", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)
		err = helpers.SeedUser(email3, "player3", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := helpers.GetUserIDByEmail(email1)
		assert.NoError(t, err)
		user2ID, err := helpers.GetUserIDByEmail(email2)
		assert.NoError(t, err)
		user3ID, err := helpers.GetUserIDByEmail(email3)
		assert.NoError(t, err)

		err = helpers.SeedSolve(user1ID, helpers.SudokusIDs[0], 30)
		assert.NoError(t, err)
		err = helpers.SeedSolve(user2ID, helpers.SudokusIDs[0], 60)
		assert.NoError(t, err)
		err = helpers.SeedSolve(user3ID, helpers.SudokusIDs[0], 45)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
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
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.Empty(t, leaderboardResp.Entries)
	})

	t.Run("get leaderboard with pagination returns correct page", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		for i := 1; i <= 5; i++ {
			email := helpers.GenerateUniqueEmail("user")
			err = helpers.SeedUser(email, "player"+string(rune('0'+i)), "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
			assert.NoError(t, err)
			userID, err := helpers.GetUserIDByEmail(email)
			assert.NoError(t, err)
			err = helpers.SeedSolve(userID, helpers.SudokusIDs[0], i*10)
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
}
