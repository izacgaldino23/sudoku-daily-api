package domain

import "sudoku-daily-api/src/domain/entities"

type (
	SudokuGenerator interface {
		GenerateDaily(size int, seed int64) *entities.Sudoku
	}
)
