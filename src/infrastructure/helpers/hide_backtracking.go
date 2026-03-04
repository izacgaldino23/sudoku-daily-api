package helpers

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
)

type (
	hideBacktracking struct {
	}
)

func NewHideBacktracking() helpers.HideBacktracking {
	return &hideBacktracking{}
}

func (s *hideBacktracking) Hide(board *entities.Sudoku, r *rand.Rand) bool {
	hideTotal := s.defineToHideCount(board, r)

	cellReference := make([][2]int, 0)
	for i := 0; i < board.GetSize(); i++ {
		for j := 0; j < board.GetSize(); j++ {
			cellReference = append(cellReference, [2]int{i, j})
		}
	}

	solver := NewSolver()

	const maxRetry = 1000

	for i := 0; i < maxRetry; i++ {
		r.Shuffle(len(cellReference), func(i, j int) {
			cellReference[i], cellReference[j] = cellReference[j], cellReference[i]
		})

		// hide numbers
		if ok := s.hideCell(board, cellReference, 0, hideTotal, solver); ok {
			for j := 0; j < hideTotal; j++ {
				board.Board.SetCell(cellReference[j][0], cellReference[j][1], 0)
			}

			return true
		}
	}

	return false
}

func (s *hideBacktracking) hideCell(board *entities.Sudoku, toHide [][2]int, current, hideTotal int, solver *Solver) bool {
	if current+1 == hideTotal {
		return true
	}

	row, col := toHide[current][0], toHide[current][1]
	n := board.Board.GetCell(row, col)

	board.Board.SetCell(row, col, 0)

	next := false
	if v := solver.Execute(board); v == 1 {
		next = s.hideCell(board, toHide, current+1, hideTotal, solver)
	}

	if !next {
		board.Board.SetCell(row, col, n)
	}

	return next
}

func (s *hideBacktracking) defineToHideCount(board *entities.Sudoku, r *rand.Rand) int {
	// get random difficulty
	min, max := entities.GetClue(board.Size, board.Difficulty)

	// get clue number between the range
	clueCount := r.Intn(max-min+1) + min
	return board.GetSize()*board.GetSize() - clueCount
}
