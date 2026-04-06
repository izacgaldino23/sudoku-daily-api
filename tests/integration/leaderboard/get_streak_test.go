package leaderboard_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/infrastructure/http/leaderboard"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func TestGetStreakLeaderboard(t *testing.T) {
	t.Run("get leaderboard with streak type returns entries", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		email := helpers.GenerateUniqueEmail("streak")
		err := helpers.SeedUser(email, "streakplayer", "$argon2id$v=19$m=65536,t=3,p=4$placeholder")
		assert.NoError(t, err)

		user1ID, err := helpers.GetUserIDByEmail(email)
		assert.NoError(t, err)

		ctx := context.Background()
		today := time.Now().Truncate(24 * time.Hour)

		_, err = database.GetDB().BunConnection.NewInsert().Model(&struct {
			bun.BaseModel  `bun:"table:user_stats"`
			ID             string    `bun:"id,pk"`
			UserID         string    `bun:"user_id"`
			CurrentStreak  int       `bun:"current_streak"`
			LongestStreak  int       `bun:"longest_streak"`
			LastSolvedDate time.Time `bun:"last_solved_date"`
			TotalSolved    int       `bun:"total_solved"`
		}{
			ID:             helpers.GenerateUUID(),
			UserID:         user1ID,
			CurrentStreak:  10,
			LongestStreak:  15,
			LastSolvedDate: today,
			TotalSolved:    50,
		}).Exec(ctx)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=streak&limit=10&page=1", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var leaderboardResp leaderboard.LeaderboardResponse
		err = json.NewDecoder(resp.Body).Decode(&leaderboardResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, leaderboardResp.Entries)
	})
}
