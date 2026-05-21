package sudoku

import (
	"net/http"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/application/usecase/sudoku"
	appContext "sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"

	"github.com/gofiber/fiber/v3"
)

type (
	Handler interface {
		GetDailySudoku(c fiber.Ctx) error
		GetDailySudokuForGuest(c fiber.Ctx) error
		CreateSudoku(c fiber.Ctx) error
		VerifySolution(c fiber.Ctx) error
		VerifySolutionGuest(c fiber.Ctx) error
		GetMyDailySudoku(c fiber.Ctx) error
		RemoveUnfinishedAttempts(c fiber.Ctx) error
	}

	sudokuHandler struct {
		getDailyUseCase                 sudoku.ISudokuGetDailyUseCase
		createSudokuUseCase             sudoku.GenerateDailyUseCase
		verifySolutionUseCase           sudoku.VerifySolutionUseCase
		verifySolutionGuestUseCase      sudoku.VerifySolutionGuestUseCase
		getUserSolvesUseCase            sudoku.GetUserSolvesUseCase
		getDailySudokuForGuest          sudoku.ISudokuGetDailyForGuestUseCase
		removeUnfinishedAttemptsUseCase sudoku.RemoveUnfinishedAttemptsUseCase
	}
)

func NewSudokuHandler(
	getDailyUseCase sudoku.ISudokuGetDailyUseCase,
	createSudokuUseCase sudoku.GenerateDailyUseCase,
	verifySolutionUseCase sudoku.VerifySolutionUseCase,
	verifySolutionGuestUseCase sudoku.VerifySolutionGuestUseCase,
	getUserSolvesUseCase sudoku.GetUserSolvesUseCase,
	getDailySudokuForGuest sudoku.ISudokuGetDailyForGuestUseCase,
	removeUnfinishedAttemptsUseCase sudoku.RemoveUnfinishedAttemptsUseCase,
) Handler {
	return &sudokuHandler{
		getDailyUseCase:                 getDailyUseCase,
		createSudokuUseCase:             createSudokuUseCase,
		verifySolutionUseCase:           verifySolutionUseCase,
		verifySolutionGuestUseCase:      verifySolutionGuestUseCase,
		getUserSolvesUseCase:            getUserSolvesUseCase,
		getDailySudokuForGuest:          getDailySudokuForGuest,
		removeUnfinishedAttemptsUseCase: removeUnfinishedAttemptsUseCase,
	}
}

// @Summary Get daily sudoku for auth users
// @Description Returns the daily sudoku puzzle for a given size for authenticated users
// @Tags sudoku
// @Accept json
// @Produce json
// @Param size query GetDailySudokuRequest true "Board size (four, six, or nine)"
// @Success 200 {object} GetDailySudokuResponse
// @Failure 400 {object} pkg.Error "invalid_query_param, invalid_size"
// @Failure 401 {object} pkg.Error "invalid_credentials, invalid_token"
// @Failure 404 {object} pkg.Error "sudoku_not_found"
// @Failure 409 {object} pkg.Error "already_played"
// @Router /api/sudoku [get]
func (sh *sudokuHandler) GetDailySudoku(c fiber.Ctx) error {
	var (
		reqCtx  = c.Context()
		request GetDailySudokuRequest
	)

	if err := c.Bind().Query(&request); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	size := entities.BoardSizeFromName(request.Size)

	userID := appContext.GetUserIDFromContext(reqCtx)
	dailySudoku, playToken, startedAt, err := sh.getDailyUseCase.Execute(reqCtx, size, userID)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response GetDailySudokuResponse
	response.FromDomain(dailySudoku, playToken, "", startedAt)

	return c.Status(http.StatusOK).JSON(response)
}

// @Summary Generate sudoku
// @Description Generates new daily sudoku puzzle for a given size
// @Tags sudoku
// @Accept json
// @Produce json
// @Param size path string true "Board size (four, six, or nine)" Enums(four, six, nine)
// @Success 200 {object} GetDailySudokuResponse
// @Failure 400 {object} pkg.Error "invalid_size"
// @Router /api/sudoku/generate/{size} [post]
func (sh *sudokuHandler) CreateSudoku(c fiber.Ctx) error {
	var (
		reqCtx = c.Context()
		size   entities.BoardSize
		err    error
	)

	sizeName := c.Params("size")
	if !entities.IsValidBoardSizeName(sizeName) {
		return pkg.JsonError(c, pkg.ErrInvalidSize)
	}

	size = entities.BoardSizeFromName(sizeName)

	dailySudoku, err := sh.createSudokuUseCase.Execute(reqCtx, size)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response GetDailySudokuResponse
	response.FromDomain(dailySudoku, "", "", time.Time{})

	return c.Status(http.StatusOK).JSON(response)
}

// @Summary Verify sudoku solution
// @Description Verifies if the submitted solution is correct
// @Tags sudoku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body VerifySolutionRequest true "Solution request"
// @Success 200 {string} string "Solution verified successfully"
// @Failure 400 {object} pkg.Error "invalid_body, invalid_solution"
// @Failure 401 {object} pkg.Error "invalid_token"
// @Failure 404 {object} pkg.Error "solution_not_found"
// @Failure 409 {object} pkg.Error "already_played"
// @Router /api/sudoku/submit [post]
func (sh *sudokuHandler) VerifySolution(c fiber.Ctx) error {
	var (
		reqCtx  = c.Context()
		now     = time.Now().UTC() // Get time here to not waste time in the use case
		err     error
		request VerifySolutionRequest
	)

	if err := c.Bind().Body(&request); err != nil {
		return pkg.ErrBodyInvalid
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	userID := appContext.GetUserIDFromContext(reqCtx)
	solve := request.ToDomain(userID)

	_, err = sh.verifySolutionUseCase.Execute(reqCtx, solve, request.PlayToken, now)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusOK)
}

// @Summary Verify guest sudoku solution
// @Description Verifies if the submitted solution is correct for guest users
// @Tags sudoku
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body VerifySolutionRequest true "Solution request"
// @Success 200 {string} string "Solution verified successfully"
// @Failure 400 {object} pkg.Error "invalid_body, invalid_solution"
// @Failure 401 {object} pkg.Error "invalid_token"
// @Router /api/sudoku/submit/guest [post]
func (sh *sudokuHandler) VerifySolutionGuest(c fiber.Ctx) error {
	var (
		reqCtx  = c.Context()
		err     error
		request VerifySolutionRequest
	)

	if err := c.Bind().Body(&request); err != nil {
		return pkg.ErrBodyInvalid
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	solve := request.ToDomain("")

	_, err = sh.verifySolutionGuestUseCase.Execute(reqCtx, solve, request.PlayToken)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusOK)
}

// @Summary Get my daily sudoku solves
// @Description Returns the daily sudoku solves for a given user
// @Tags sudoku
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MySolvesResponse
// @Failure 401 {object} pkg.Error "invalid_token"
// @Router /api/sudoku/me [get]
func (sh *sudokuHandler) GetMyDailySudoku(c fiber.Ctx) error {
	var (
		reqCtx = c.Context()
		err    error
	)

	userID := appContext.GetUserIDFromContext(reqCtx)
	solves, err := sh.getUserSolvesUseCase.Execute(reqCtx, userID)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response MySolvesResponse
	response.FromDomain(solves)

	return c.Status(http.StatusOK).JSON(response)
}

// @Summary Get daily sudoku
// @Description Returns the daily sudoku puzzle for a given size
// @Tags sudoku
// @Accept json
// @Produce json
// @Param size query GetDailySudokuRequest true "Board size (four, six, or nine)"
// @Success 200 {object} GetDailySudokuResponse
// @Failure 400 {object} pkg.Error "invalid_query_param, invalid_size"
// @Failure 404 {object} pkg.Error "sudoku_not_found"
// @Router /api/sudoku [get]
func (sh *sudokuHandler) GetDailySudokuForGuest(c fiber.Ctx) error {
	var (
		reqCtx  = c.Context()
		request GetDailySudokuRequest
	)

	if err := c.Bind().Query(&request); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	size := entities.BoardSizeFromName(request.Size)

	dailySudoku, playToken, sessionID, err := sh.getDailySudokuForGuest.Execute(reqCtx, size)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response GetDailySudokuResponse
	response.FromDomain(dailySudoku, playToken, sessionID, time.Time{})

	return c.Status(http.StatusOK).JSON(response)
}

// @Summary Remove unfinished attempts
// @Description Removes unfinished attempts for the daily sudoku puzzles that are past the reset threshold (e.g., 24 hours). This endpoint is intended to be called by a scheduled job to clean up old attempts and reset strikes for users who haven't completed their puzzles in time.
// @Tags sudoku
// @Accept json
// @Produce json
// @Success 200 "Cleaned unfinished attempts"
// @Router /api/cron/unfinished-attempts [post]
func (sh *sudokuHandler) RemoveUnfinishedAttempts(c fiber.Ctx) error {
	reqCtx := c.Context()

	err := sh.removeUnfinishedAttemptsUseCase.Execute(reqCtx)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusOK)
}
