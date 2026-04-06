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
		ResetStrikes(c fiber.Ctx) error
	}

	leaderboardHandler struct {
		leaderboardUsecase  usecase.GetLeaderboard
		resetStrikesUseCase usecase.ResetStrikesUseCase
	}
)

func NewLeaderboardHandler(leaderboardUsecase usecase.GetLeaderboard, resetStrikesUseCase usecase.ResetStrikesUseCase) LeaderboardHandler {
	return &leaderboardHandler{
		leaderboardUsecase:  leaderboardUsecase,
		resetStrikesUseCase: resetStrikesUseCase,
	}
}

// @Summary Get leaderboard
// @Description Returns the leaderboard with rankings for a given type and size. For daily and all-time types, size is required. For streak and total types, size should not be provided.
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param type query string true "Leaderboard type (daily, all-time, streak, total)"
// @Param size query string false "Board size (four, six, nine) - required for daily and all-time, not allowed for streak and total"
// @Param limit query int false "Number of entries to return (1-100)"
// @Param page query int false "Page number"
// @Success 200 {object} LeaderboardResponse
// @Failure 400 {object} pkg.Error "invalid_leaderboard_type, invalid_size, invalid_limit, invalid_page, size_required, size_not_allowed"
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

	if err = params.Validate(); err != nil {
		return pkg.JsonError(c, err)
	}

	leaderboard, err := h.leaderboardUsecase.Execute(reqCtx, params.ToDomain())
	if err != nil {
		return pkg.JsonError(c, err)
	}

	response := responseFromDomain(leaderboard)

	return c.
		Status(http.StatusOK).
		JSON(response)
}

func (h *leaderboardHandler) ResetStrikes(c fiber.Ctx) error {
	var (
		req    ResetStrikesRequest
		err    error
		reqCtx = c.Context()
	)

	if err = c.Bind().Body(&req); err != nil {
		return pkg.ErrBodyInvalid
	}

	if err = pkg.ValidateStruct(req); err != nil {
		return pkg.JsonError(c, err)
	}

	err = h.resetStrikesUseCase.Execute(reqCtx, req.Date)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusNoContent)
}
