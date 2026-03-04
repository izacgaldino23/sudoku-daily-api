package domain

import "sudoku-daily-api/src/domain/entities"

type (
	SudokuGenerator interface {
		GenerateDaily(size entities.BoardSize, seed int64) (*entities.Sudoku, error)
	}
)
