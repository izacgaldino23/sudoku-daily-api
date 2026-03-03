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

func (h *hideBacktracking) Hide(board *entities.Sudoku, r *rand.Rand) {
	h.hideNumbers(board, r)
}

func (s *hideBacktracking) hideNumbers(board *entities.Sudoku, r *rand.Rand) {
	hideTotal := s.defineToHideCount(board, r)

	cellReference := make([][2]int, 0)
	for i := 0; i < board.GetSize(); i++ {
		for j := 0; j < board.GetSize(); j++ {
			cellReference = append(cellReference, [2]int{i, j})
		}
	}

	solver := NewSolver()
	current := 0

	for {
		if current + hideTotal > len(cellReference) {
			current = 0
		}

		if current == 0 {
			r.Shuffle(len(cellReference), func(i, j int) {
				cellReference[i], cellReference[j] = cellReference[j], cellReference[i]
			})
		}

		// hide numbers
		if ok := s.hideCell(board, cellReference, current, hideTotal, solver); ok {
			break
		}

		current++
	}
}

func (s *hideBacktracking) hideCell(board *entities.Sudoku, toHide [][2]int, current, hideTotal int, solver *Solver) bool {
	if current == hideTotal {
		return true
	}

	row, col := toHide[current][0], toHide[current][1]
	n := board.Board.GetCell(row, col)

	board.Board.SetCell(row, col, 0)

	// test solutions
	solutions := solver.Execute(board)
	if solutions != 1 {
		board.Board.SetCell(row, col, n)
		return false
	}

	return s.hideCell(board, toHide, current+1, hideTotal, solver)
}

func (s *hideBacktracking) defineToHideCount(board *entities.Sudoku, r *rand.Rand) int {
	difficulties := []entities.Difficulty{
		entities.DifficultyEasy,
		entities.DifficultyMedium,
		entities.DifficultyHard,
	}

	// get random difficulty
	difficulty := difficulties[r.Intn(len(difficulties))]
	min, max := entities.GetClue(entities.BoardSize(board.Size), difficulty)

	// get clue number between the range
	clueCount := r.Intn(max-min+1) + min
	return board.GetSize()*board.GetSize() - clueCount
}
