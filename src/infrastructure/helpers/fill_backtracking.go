package helpers

import (
	"math/rand"
	"slices"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
)

type (
	fillBacktracking struct {
		baseNumber []int
	}
)

func NewFillBacktracking() helpers.FillBacktracking {
	return &fillBacktracking{}
}

func (f *fillBacktracking) Fill(board *entities.Sudoku, r *rand.Rand) {
	f.baseNumber = make([]int, board.Size)
	for i := 0; i < board.Size; i++ {
		f.baseNumber[i] = i + 1
	}

	f.fillCell(board, 0, 0, []int{}, r, board.GetGrids(board.Size))
}

func (f *fillBacktracking) fillCell(board *entities.Sudoku, currentRow, currentCol int, chosen []int, r *rand.Rand, grids []entities.Grid) bool {
	var currentDecision []int

	if len(chosen) == 0 {
		currentDecision = make([]int, board.Size-len(chosen))
		copy(currentDecision, f.baseNumber)
	} else {
		for i := range board.Size {
			if !slices.Contains(chosen, i+1) {
				currentDecision = append(currentDecision, i+1)
			}
		}
	}

	// shuffle numbers
	if len(currentDecision) > 1 {
		r.Shuffle(len(currentDecision), func(i, j int) {
			currentDecision[i], currentDecision[j] = currentDecision[j], currentDecision[i]
		})
	}

	for i := range currentDecision {
		n := currentDecision[i]
		board.Board[currentRow][currentCol] = n

		// validate line
		if !isLineValid(board.Board, currentRow, 0, 1, board.Size) {
			continue
		}

		// validate columns
		if !isLineValid(board.Board, 0, currentCol, board.Size, 1) {
			continue
		}

		// validate grid
		valid := true
		for _, grid := range grids {
			if grid.IsLastPosition(currentRow, currentCol) {
				if !isLineValid(board.Board, grid.Row, grid.Col, grid.RowCount, grid.ColCount) {
					valid = false
				}
			}
		}
		if !valid {
			continue
		}

		if currentCol == board.Size-1 && currentRow == board.Size-1 {
			return true
		}

		if currentCol == board.Size-1 {
			if f.fillCell(board, currentRow+1, 0, []int{}, r, grids) {
				return true
			}
		} else {
			// call the next cell
			if f.fillCell(board, currentRow, currentCol+1, append(chosen, n), r, grids) {
				return true
			}
		}
	}

	board.Board[currentRow][currentCol] = 0

	return false
}
