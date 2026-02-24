package sudoku

import (
	"net/http"
	"strconv"
	"sudoku-daily-api/src/adapters/helpers"
	"sudoku-daily-api/src/application/usecase"
	"sudoku-daily-api/src/domain/entities"

	"github.com/gofiber/fiber/v3"
)

type (
	ISudokuHandler interface {
		GetDailySudoku(c fiber.Ctx) error
	}

	sudokuHandler struct {
		getDailyUseCase usecase.ISudokuGetDailyUseCase
	}
)

func NewSudokuHandler(getDailyUseCase usecase.ISudokuGetDailyUseCase) ISudokuHandler {
	return &sudokuHandler{getDailyUseCase: getDailyUseCase}
}

func (sh *sudokuHandler) GetDailySudoku(c fiber.Ctx) error {
	var (
		sizeParam = c.Params("size")
		ctxReq    = c.Context()
		size      int
		err       error
	)

	if sizeParam == "" {
		return helpers.JsonError(c, "Invalid size")
	}

	size, err = strconv.Atoi(sizeParam)
	if err != nil {
		return helpers.JsonError(c, "Invalid size")
	}

	var dailySudoku *entities.Sudoku
	dailySudoku, err = sh.getDailyUseCase.Execute(ctxReq, size)
	if err != nil {
		return helpers.JsonErrorWithStatus(c, err.Error(), http.StatusInternalServerError)
	}

	var response GetDailySudokuResponse
	response.FromDomain(dailySudoku)

	return c.Status(http.StatusOK).JSON(response)
}
