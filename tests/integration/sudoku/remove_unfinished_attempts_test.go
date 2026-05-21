package sudoku_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

type solveRow struct {
	bun.BaseModel `bun:"table:solves"`
	ID            string    `bun:"id,pk"`
	UserID        string    `bun:"user_id"`
	SudokuID      string    `bun:"sudoku_id"`
	StartedAt     time.Time `bun:"started_at"`
	Duration      int       `bun:"duration"`
	Size          int       `bun:"size"`
}

func countSolves() (int, error) {
	db := database.GetDB().BunConnection
	ctx := context.Background()
	count, err := db.NewSelect().Model(&solveRow{}).Count(ctx)
	return count, err
}

func TestRemoveUnfinishedAttempts(t *testing.T) {
	db := database.GetDB().BunConnection
	ctx := context.Background()
	today := time.Now().Truncate(24 * time.Hour)

	t.Run("removes unfinished attempts within window", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		localUser, err := helpers.RegisterAndLoginUser(app, "strong-password-123")
		assert.NoError(t, err)
		localID, err := helpers.GetUserIDByEmail(localUser.Email)
		assert.NoError(t, err)

		helpers.SeedSolve(localID, helpers.SudokusIDs[0], 0)
		db.NewUpdate().Model(&solveRow{}).Set("started_at = ?", today.Add(-24*time.Hour)).Where("user_id = ?", localID).Exec(ctx)

		req := httptest.NewRequest(http.MethodPost, "/api/cron/unfinished-attempts", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		count, err := countSolves()
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("preserves completed attempts", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		localUser, err := helpers.RegisterAndLoginUser(app, "strong-password-123")
		assert.NoError(t, err)
		localID, err := helpers.GetUserIDByEmail(localUser.Email)
		assert.NoError(t, err)

		helpers.SeedSolve(localID, helpers.SudokusIDs[0], 120)

		req := httptest.NewRequest(http.MethodPost, "/api/cron/unfinished-attempts", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		count, err := countSolves()
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("preserves attempts outside date range", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		localUser, err := helpers.RegisterAndLoginUser(app, "strong-password-123")
		assert.NoError(t, err)
		localID, err := helpers.GetUserIDByEmail(localUser.Email)
		assert.NoError(t, err)

		helpers.SeedSolve(localID, helpers.SudokusIDs[0], 0)
		db.NewUpdate().Model(&solveRow{}).Set("started_at = ?", today.Add(-72*time.Hour)).Where("user_id = ?", localID).Exec(ctx)

		req := httptest.NewRequest(http.MethodPost, "/api/cron/unfinished-attempts", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		count, err := countSolves()
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("returns 200 when no attempts exist", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodPost, "/api/cron/unfinished-attempts", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("only removes matching unfinished attempts", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		localUser, err := helpers.RegisterAndLoginUser(app, "strong-password-123")
		assert.NoError(t, err)
		localID, err := helpers.GetUserIDByEmail(localUser.Email)
		assert.NoError(t, err)

		helpers.SeedSolve(localID, helpers.SudokusIDs[0], 0)
		db.NewUpdate().Model(&solveRow{}).Set("started_at = ?", today.Add(-24*time.Hour)).Where("user_id = ?", localID).Where("sudoku_id = ?", helpers.SudokusIDs[0]).Exec(ctx)

		helpers.SeedSolve(localID, helpers.SudokusIDs[1], 60)
		helpers.SeedSolve(localID, helpers.SudokusIDs[2], 0)
		db.NewUpdate().Model(&solveRow{}).Set("started_at = ?", today.Add(-72*time.Hour)).Where("user_id = ?", localID).Where("sudoku_id = ?", helpers.SudokusIDs[2]).Exec(ctx)

		req := httptest.NewRequest(http.MethodPost, "/api/cron/unfinished-attempts", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		count, err := countSolves()
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("idempotent", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		err := helpers.SeedSudokus()
		assert.NoError(t, err)

		localUser, err := helpers.RegisterAndLoginUser(app, "strong-password-123")
		assert.NoError(t, err)
		localID, err := helpers.GetUserIDByEmail(localUser.Email)
		assert.NoError(t, err)

		helpers.SeedSolve(localID, helpers.SudokusIDs[0], 0)
		db.NewUpdate().Model(&solveRow{}).Set("started_at = ?", today.Add(-24*time.Hour)).Where("user_id = ?", localID).Exec(ctx)

		req := httptest.NewRequest(http.MethodPost, "/api/cron/unfinished-attempts", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		count, err := countSolves()
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
