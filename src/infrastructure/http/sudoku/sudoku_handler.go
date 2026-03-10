package sudoku

import (
	"net/http"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/application/usecase/sudoku"
	appContext "sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"

	"github.com/gofiber/fiber/v3"
)

type (
	ISudokuHandler interface {
		GetDailySudoku(c fiber.Ctx) error
		CreateSudoku(c fiber.Ctx) error
		VerifySolution(c fiber.Ctx) error
	}

	sudokuHandler struct {
		getDailyUseCase       sudoku.ISudokuGetDailyUseCase
		createSudokuUseCase   sudoku.SudokuGenerateAllUseCase
		verifySolutionUseCase sudoku.SudokuVerifySolutionUseCase
	}
)

func NewSudokuHandler(
	getDailyUseCase sudoku.ISudokuGetDailyUseCase,
	createSudokuUseCase sudoku.SudokuGenerateAllUseCase,
	verifySolutionUseCase sudoku.SudokuVerifySolutionUseCase,
) ISudokuHandler {
	return &sudokuHandler{
		getDailyUseCase:       getDailyUseCase,
		createSudokuUseCase:   createSudokuUseCase,
		verifySolutionUseCase: verifySolutionUseCase,
	}
}

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

	size := request.GetSize()

	dailySudoku, playToken, sessionID, err := sh.getDailyUseCase.Execute(ctxReq, size)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	var response SudokuResponse
	response.FromDomain(dailySudoku, playToken, sessionID)

	return c.Status(http.StatusOK).JSON(response)
}

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
		s.FromDomain(&sudoku, "", "")
		response = append(response, s)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (sh *sudokuHandler) VerifySolution(c fiber.Ctx) error {
	var (
		ctxReq  = c.Context()
		err     error
		request VerifySolutionRequest
	)

	if err := c.Bind().Body(&request); err != nil {
		return pkg.JsonErrorWithStatus(c, err, http.StatusBadRequest)
	}

	if err := pkg.ValidateStruct(request); err != nil {
		return pkg.JsonError(c, err)
	}

	userID := appContext.GetUserIDFromContext(ctxReq)
	solve := request.ToDomain(userID)

	_, err = sh.verifySolutionUseCase.Execute(ctxReq, solve, request.PlayToken)
	if err != nil {
		return pkg.JsonError(c, err)
	}

	return c.SendStatus(http.StatusOK)
}
