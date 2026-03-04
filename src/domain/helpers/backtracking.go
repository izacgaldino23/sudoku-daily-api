package helpers

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
)

type (
	FillBacktracking interface {
		Fill(board *entities.Sudoku, r *rand.Rand) bool
	}

	HideBacktracking interface {
		Hide(board *entities.Sudoku, r *rand.Rand) bool
	}
)
