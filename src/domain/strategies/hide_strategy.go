package strategies

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
)

type (
	HideStrategy interface {
		Hide(board *entities.Sudoku, r *rand.Rand) bool
	}

	hideBacktracking struct {
		solver *solver
	}
)

func NewHideStrategy() HideStrategy {
	return &hideBacktracking{
		solver: newSolver(),
	}
}

func (s *hideBacktracking) Hide(board *entities.Sudoku, r *rand.Rand) bool {
	targetToHide := s.defineToHideCount(board, r)

	const maxTries = 1000

	for i := 0; i < maxTries; i++ {
		cells := s.getCellShuffled(board, r)
		var hidden int
		for j := range cells {
			cell := cells[j]
			if hidden >= targetToHide {
				return true
			}

			val := board.Board.GetCell(cell[0], cell[1])
			board.Board.SetCell(cell[0], cell[1], 0)

			if s.solver.Execute(board) == 1 {
				hidden++
			} else {
				board.Board.SetCell(cell[0], cell[1], val)
			}
		}
	}

	return false
}

func (s *hideBacktracking) defineToHideCount(board *entities.Sudoku, r *rand.Rand) int {
	min, max := entities.GetClue(board.Size, board.Difficulty)

	clueCount := r.Intn(max-min+1) + min
	return board.GetSize()*board.GetSize() - clueCount
}

func (s *hideBacktracking) getCellShuffled(board *entities.Sudoku, r *rand.Rand) [][2]int {
	cellReference := make([][2]int, 0)

	for row := range board.Board.GetBoard() {
		for col := range board.Board.GetBoard()[row] {
			if board.Board.GetCell(row, col) != 0 {
				cellReference = append(cellReference, [2]int{row, col})
			}
		}
	}

	r.Shuffle(len(cellReference), func(i, j int) {
		cellReference[i], cellReference[j] = cellReference[j], cellReference[i]
	})

	return cellReference
}
