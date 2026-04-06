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
	SudokuHandler interface {
		GetDailySudoku(c fiber.Ctx) error
		CreateSudoku(c fiber.Ctx) error
		VerifySolution(c fiber.Ctx) error
		GetMyDailySudoku(c fiber.Ctx) error
	}

	sudokuHandler struct {
		getDailyUseCase       sudoku.ISudokuGetDailyUseCase
		createSudokuUseCase   sudoku.SudokuGenerateDailyUseCase
		verifySolutionUseCase sudoku.SudokuVerifySolutionUseCase
		getUserSolvesUseCase  sudoku.SudokuGetUserSolvesUseCase
	}
)

func NewSudokuHandler(
	getDailyUseCase sudoku.ISudokuGetDailyUseCase,
	createSudokuUseCase sudoku.SudokuGenerateDailyUseCase,
	verifySolutionUseCase sudoku.SudokuVerifySolutionUseCase,
	getUserSolvesUseCase sudoku.SudokuGetUserSolvesUseCase,
) SudokuHandler {
	return &sudokuHandler{
		getDailyUseCase:       getDailyUseCase,
		createSudokuUseCase:   createSudokuUseCase,
		verifySolutionUseCase: verifySolutionUseCase,
		getUserSolvesUseCase:  getUserSolvesUseCase,
	}
}

// @Summary Get daily sudoku
// @Description Returns the daily sudoku puzzle for a given size
// @Tags sudoku
// @Accept json
// @Produce json
// @Param size query GetDailySudokuRequest true "Board size (four, six, or nine)"
// @Success 200 {object} SudokuResponse
// @Failure 400 {object} pkg.Error
// @Failure 404 {object} pkg.Error
// @Failure 409 {object} pkg.Error
// @Router /api/sudoku [get]
func (sh *sudokuHandler) GetDailySudoku(c fiber.Ctx) error {
	var (
		ctxReq  = c.Context()
		request GetDailySudokuRequest
	)

	if err := c.Bind().Query(&request); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	size := entities.BoardSizeFromName(request.Size)

	userID := appContext.GetUserIDFromContext(ctxReq)
	dailySudoku, playToken, err := sh.getDailyUseCase.Execute(ctxReq, size, userID)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response SudokuResponse
	response.FromDomain(dailySudoku, playToken)

	return c.Status(http.StatusOK).JSON(response)
}

// @Summary Generate sudoku
// @Description Generates new daily sudoku puzzles for all sizes
// @Tags sudoku
// @Produce json
// @Success 200 {array} SudokuResponse
// @Router /api/sudoku/generate [post]
func (sh *sudokuHandler) CreateSudoku(c fiber.Ctx) error {
	var (
		ctxReq = c.Context()
		err    error
	)

	var dailySudoku []entities.Sudoku
	dailySudoku, err = sh.createSudokuUseCase.Execute(ctxReq)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response []SudokuResponse
	for _, sudoku := range dailySudoku {
		s := SudokuResponse{}
		s.FromDomain(&sudoku, "")
		response = append(response, s)
	}

	return c.Status(http.StatusOK).JSON(response)
}

// @Summary Verify sudoku solution
// @Description Verifies if the submitted solution is correct
// @Tags sudoku
// @Accept json
// @Produce json
// @Security BearerAuth  // optional
// @Param request body VerifySolutionRequest true "Solution request"
// @Success 200 {string} string "Solution verified successfully"
// @Failure 400 {object} pkg.Error
// @Failure 404 {object} pkg.Error
// @Failure 409 {object} pkg.Error
// @Router /api/sudoku/submit [post]
func (sh *sudokuHandler) VerifySolution(c fiber.Ctx) error {
	var (
		ctxReq  = c.Context()
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

	userID := appContext.GetUserIDFromContext(ctxReq)
	solve := request.ToDomain(userID)

	_, err = sh.verifySolutionUseCase.Execute(ctxReq, solve, request.PlayToken, now)
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
// @Failure 401 {object} pkg.Error
// @Router /api/sudoku/me [get]
func (sh *sudokuHandler) GetMyDailySudoku(c fiber.Ctx) error {
	var (
		ctxReq = c.Context()
		err    error
	)

	userID := appContext.GetUserIDFromContext(ctxReq)
	solves, err := sh.getUserSolvesUseCase.Execute(ctxReq, userID)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response MySolvesResponse
	response.FromDomain(solves)

	return c.Status(http.StatusOK).JSON(response)
}
