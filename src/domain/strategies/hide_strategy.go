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

		if s.solveRecursive(&board.Board, cells, 0, 0, targetToHide) {
			return true
		}
	}

	return false
}

func (s *hideBacktracking) solveRecursive(board *entities.Board, cells [][2]int, index int, hiddenCount int, target int) bool {
    if hiddenCount >= target {
        return true
    }

    if index >= len(cells) {
        return false
    }

    cell := cells[index]
    row, col := cell[0], cell[1]
    originalVal := board.GetCell(row, col)

    board.SetCell(row, col, 0)

    if s.solver.Execute(board) == 1 {
        if s.solveRecursive(board, cells, index+1, hiddenCount+1, target) {
            return true
        }
    }

    board.SetCell(row, col, originalVal)

    // Tenta esconder as próximas SEM esconder esta atual
    return s.solveRecursive(board, cells, index+1, hiddenCount, target)
}

func (s *hideBacktracking) defineToHideCount(board *entities.Sudoku, r *rand.Rand) int {
	if board.Difficulty == "" {
		// generate random difficulty
		difficulties := []entities.Difficulty{entities.DifficultyEasy, entities.DifficultyMedium, entities.DifficultyHard}
		board.Difficulty = difficulties[r.Intn(len(difficulties))]
	}

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
