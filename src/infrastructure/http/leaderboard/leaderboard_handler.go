package leaderboard

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"sudoku-daily-api/pkg"
	usecase "sudoku-daily-api/src/application/usecase/leaderboard"
)

type (
	LeaderboardHandler interface {
		GetLeaderboard(c fiber.Ctx) error
	}

	leaderboardHandler struct {
		leaderboardUsecase usecase.GetLeaderboard
	}
)

func NewLeaderboardHandler(leaderboardUsecase usecase.GetLeaderboard) LeaderboardHandler {
	return &leaderboardHandler{
		leaderboardUsecase: leaderboardUsecase,
	}
}

// @Summary Get leaderboard
// @Description Returns the leaderboard with rankings for a given type and size
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param type query LeaderboardRequest false "Leaderboard type (daily, all-time, streak, total)"
// @Param size query LeaderboardRequest false "Board size (four, six, nine)"
// @Param limit query int false "Number of entries to return (1-100)"
// @Param page query int false "Page number"
// @Success 200 {object} LeaderboardResponse
// @Failure 400 {object} pkg.Error
// @Router /api/leaderboard [get]
func (h *leaderboardHandler) GetLeaderboard(c fiber.Ctx) error {
	var (
		params LeaderboardRequest
		err    error
		reqCtx = c.Context()
	)

	if err = c.Bind().Query(&params); err != nil {
		return pkg.JsonError(c, pkg.ErrQueryParamInvalid)
	}

	if err = pkg.ValidateStruct(params); err != nil {
		return pkg.JsonError(c, err)
	}

	leaderboard, err := h.leaderboardUsecase.Execute(reqCtx, params.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.
		Status(http.StatusOK).
		JSON(responseFromDomain(leaderboard))
}
