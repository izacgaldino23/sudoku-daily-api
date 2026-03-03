package helpers

import (
	"fmt"
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
)

type (
	fillBacktracking struct{}
)

func NewFillBacktracking() helpers.FillBacktracking {
	return &fillBacktracking{}
}

func (f *fillBacktracking) Fill(board *entities.Sudoku, r *rand.Rand) {
	f.fillCell(board, 0, 0, r)
}

func (f *fillBacktracking) fillCell(board *entities.Sudoku, currentRow, currentCol int, r *rand.Rand) bool {
	if currentRow == board.GetSize() {
		return true
	}

	missing := board.Board.GetPossibleByPosition(currentRow, currentCol)
	values := missing.Values()

	r.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	for _, n := range missing.Values() {
		if board.Board.HasNumber(currentRow, currentCol, n) {
			continue
		}

		board.Board.SetCell(currentRow, currentCol, n)
		fmt.Printf("+ Row: %v, Col: %v, Value: %v, size: %v\n", currentRow, currentCol, n, board.Size)

		// go to next in the same row
		if currentCol == board.GetSize()-1 {
			if f.fillCell(board, currentRow+1, 0, r) {
				return true
			}
		} else {
			// call the next row
			if f.fillCell(board, currentRow, currentCol+1, r) {
				return true
			}
		}

		// backtracking
		fmt.Printf("- Row: %v, Col: %v, Value: %v, size: %v\n", currentRow, currentCol, n, board.Size)
		board.Board.SetCell(currentRow, currentCol, 0)
	}

	return false
}
