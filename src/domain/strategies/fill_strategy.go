package strategies

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
)

type (
	FillStrategy interface {
		Fill(board *entities.Sudoku, r *rand.Rand) bool
	}

	fillBacktracking struct{}
)

func NewFillStrategy() FillStrategy {
	return &fillBacktracking{}
}

func (f *fillBacktracking) Fill(board *entities.Sudoku, r *rand.Rand) bool {
	return f.fillCell(&board.Board, 0, 0, r)
}

func (f *fillBacktracking) fillCell(board *entities.Board, currentRow, currentCol int, r *rand.Rand) bool {
	if currentRow == board.GetSize() {
		return true
	}

	missing := board.GetPossibleByPosition(currentRow, currentCol)
	values := missing.Values()

	r.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	for _, n := range missing.Values() {
		if board.HasNumber(currentRow, currentCol, n) {
			continue
		}

		board.SetCell(currentRow, currentCol, n)

		if currentCol == board.GetSize()-1 {
			if f.fillCell(board, currentRow+1, 0, r) {
				return true
			}
		} else {
			if f.fillCell(board, currentRow, currentCol+1, r) {
				return true
			}
		}

		board.SetCell(currentRow, currentCol, 0)
	}

	return false
}
