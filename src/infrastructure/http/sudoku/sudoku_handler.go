package sudoku

import (
	"net/http"
	"strconv"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/application/usecase"
	"sudoku-daily-api/src/domain/entities"

	"github.com/gofiber/fiber/v3"
)

type (
	ISudokuHandler interface {
		GetDailySudoku(c fiber.Ctx) error
		CreateSudoku(c fiber.Ctx) error
	}

	sudokuHandler struct {
		getDailyUseCase     usecase.ISudokuGetDailyUseCase
		createSudokuUseCase usecase.ISudokuGenerateAllUseCase
	}
)

func NewSudokuHandler(
	getDailyUseCase usecase.ISudokuGetDailyUseCase,
	createSudokuUseCase usecase.ISudokuGenerateAllUseCase,
) ISudokuHandler {
	return &sudokuHandler{
		getDailyUseCase:     getDailyUseCase,
		createSudokuUseCase: createSudokuUseCase,
	}
}

func (sh *sudokuHandler) GetDailySudoku(c fiber.Ctx) error {
	var (
		sizeParam = c.Params("size")
		ctxReq    = c.Context()
		size      int
		err       error
	)

	if sizeParam == "" {
		return pkg.JsonError(c, "Invalid size")
	}

	size, err = strconv.Atoi(sizeParam)
	if err != nil {
		return pkg.JsonError(c, "Invalid size")
	}

	var dailySudoku *entities.Sudoku
	dailySudoku, err = sh.getDailyUseCase.Execute(ctxReq, size)
	if err != nil {
		return pkg.JsonErrorWithStatus(c, err.Error(), http.StatusInternalServerError)
	}

	var response GetDailySudokuResponse
	response.FromDomain(dailySudoku)

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
		return pkg.JsonErrorWithStatus(c, err.Error(), http.StatusInternalServerError)
	}

	var response []GetDailySudokuResponse
	for _, sudoku := range dailySudoku {
		s := GetDailySudokuResponse{}
		s.FromDomain(&sudoku)
		response = append(response, s)
	}

	return c.Status(http.StatusOK).JSON(response)
}
