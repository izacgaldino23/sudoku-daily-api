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

func TestGetAllTimeLeaderboard(t *testing.T) {
	t.Run("get leaderboard with all-time type returns entries", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		email := helpers.GenerateUniqueEmail("user1")
		err = helpers.SeedUser(email, "bestplayer", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := helpers.GetUserIDByEmail(email)
		assert.NoError(t, err)

		err = helpers.SeedSolves(user1ID)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=all-time&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardResp.Entries)
	})
}
