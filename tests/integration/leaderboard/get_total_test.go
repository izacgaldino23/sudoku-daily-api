package leaderboard_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/infrastructure/http/leaderboard"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalLeaderboard(t *testing.T) {
	t.Run("get leaderboard with total type returns entries", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		email := helpers.GenerateUniqueEmail("user1")
		err = helpers.SeedUser(email, "activeplayer", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := helpers.GetUserIDByEmail(email)
		assert.NoError(t, err)

		err = helpers.SeedSolves(user1ID)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=total&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardResp.Entries)
	})

	t.Run("get leaderboard with total type and size returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=total&size=nine&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var errorResp pkg.Error
		err = json.Unmarshal(body, &errorResp)
		assert.NoError(t, err)
		assert.Len(t, errorResp.ValidationErr, 1)
		assert.Equal(t, "Size", errorResp.ValidationErr[0].Field)
	})
}
