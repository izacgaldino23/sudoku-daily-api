package domain

import "sudoku-daily-api/src/domain/entities"

type (
	Generator interface {
		GenerateDaily(size int, difficulty string, seed int64) *entities.Sudoku
	}
)
