package helpers

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
	"sudoku-daily-api/src/domain/vo"
)

type (
	fillBacktracking struct {
		baseNumbers vo.Binary
	}
)

func NewFillBacktracking() helpers.FillBacktracking {
	return &fillBacktracking{}
}

func (f *fillBacktracking) Fill(board *entities.Sudoku, r *rand.Rand) {
	for i := 0; i < board.Size; i++ {
		f.baseNumbers.Add(i + 1)
	}

	f.fillCell(board, 0, 0, 0, r, board.GetGrids(board.Size))
}

func (f *fillBacktracking) fillCell(board *entities.Sudoku, currentRow, currentCol int, chosen vo.Binary, r *rand.Rand, grids []entities.Grid) bool {
	var currentDecision = []int{}

	for i := range board.Size {
		if !chosen.Contains(i + 1) {
			currentDecision = append(currentDecision, i+1)
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
		gridRows := entities.BoardSizes[entities.BoardSize(board.Size)]
		gridCols := board.Size / gridRows
		if currentRow+1%gridRows == 0 {
			gridCol := (currentCol / gridCols) * gridCols
			if !isLineValid(board.Board, currentRow-gridRows+1, gridCol, gridRows, gridCols) {
				valid = false
			}
		}
		if !valid {
			continue
		}

		if currentCol == board.Size-1 && currentRow == board.Size-1 {
			return true
		}

		if currentCol == board.Size-1 {
			if f.fillCell(board, currentRow+1, 0, 0, r, grids) {
				return true
			}
		} else {
			// call the next cell
			chosen.Add(n)
			if f.fillCell(board, currentRow, currentCol+1, chosen, r, grids) {
				return true
			}
			chosen.Remove(n)
		}
	}

	board.Board[currentRow][currentCol] = 0

	return false
}
