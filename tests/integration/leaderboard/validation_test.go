package leaderboard_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sudoku-daily-api/tests/integration/helpers"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardValidation(t *testing.T) {
	t.Run("get leaderboard with invalid type returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=invalid&size=nine", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with invalid size returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=invalid", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with missing type returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?size=nine", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with missing size returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with limit exceeding max returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=500", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with limit below min returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&limit=0", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("get leaderboard with page below min returns bad request", func(t *testing.T) {
		t.Cleanup(helpers.TruncateTables)
		app := helpers.SetupTestApp()

		req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily&size=nine&page=0", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
