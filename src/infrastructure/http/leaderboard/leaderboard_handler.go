package leaderboard

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"sudoku-daily-api/pkg"
	usecase "sudoku-daily-api/src/application/usecase/leaderboard"
	"sudoku-daily-api/src/domain/entities"
)

type (
	LeaderboardHandler interface {
		GetLeaderboard(c fiber.Ctx) error
	}

	leaderboardHandler struct {
		leaderboardUsecase usecase.LeaderboardUsecase
	}
)

func NewLeaderboardHandler(leaderboardUsecase usecase.LeaderboardUsecase) LeaderboardHandler {
	return &leaderboardHandler{
		leaderboardUsecase: leaderboardUsecase,
	}
}

func (h *leaderboardHandler) GetLeaderboard(c fiber.Ctx) error {
	var (
		params entities.LeaderboardSearchParams
		err    error
		reqCtx = c.Context()
	)

	if err = c.Bind().Query(&params); err != nil {
		return pkg.JsonError(c, pkg.ErrQueryParamInvalid)
	}

	if err = pkg.ValidateStruct(params); err != nil {
		return pkg.JsonError(c, err)
	}

	leaderboard, err := h.leaderboardUsecase.Execute(reqCtx, params)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.
		Status(http.StatusOK).
		JSON(responseFromDomain(&leaderboard))
}
