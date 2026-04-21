package strategies

import (
	"context"
	"math/rand"
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/infrastructure/logging"
)

type (
	HideStrategy interface {
		Hide(ctx context.Context, board *entities.Sudoku, r *rand.Rand) bool
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

func (s *hideBacktracking) Hide(ctx context.Context, board *entities.Sudoku, r *rand.Rand) bool {
	targetToHide := s.defineToHideCount(board, r)

	const (
		maxTries          = 1000
		maxTargetDecrease = 15
	)

	var (
		tries     int
		startTime = time.Now()
	)

	for range maxTargetDecrease {
		for range maxTries {
			cells := s.getCellShuffled(board, r)

			if s.hideCells(&board.Board, cells, targetToHide) {
				logging.Log(ctx).Info().Msgf("Successfully hidden %v cells, tries: %v, board size: %v, time: %v", targetToHide, tries, board.Size, time.Since(startTime))
				return true
			}

			tries++
		}

		targetToHide--
	}

	logging.Log(ctx).Error().Msgf("Failed to hide %v cells, tries: %v, board size: %v", targetToHide, tries, board.Size)

	return false
}

func (s *hideBacktracking) hideCells(board *entities.Board, cells [][3]int, target int) bool {
	var (
		hiddenCount = 0
	)
	for i := range cells {
		if hiddenCount >= target || i+(target-hiddenCount) >= len(cells) {
			break
		}

		cell := cells[i]
		row, col, value := cell[0], cell[1], cell[2]

		board.SetCell(row, col, 0)
		if s.solver.Execute(board) == 1 {
			hiddenCount++
		} else {
			board.SetCell(row, col, value)
		}
	}

	hidedTarget := hiddenCount >= target

	if !hidedTarget {
		for i := range cells {
			cell := cells[i]
			row, col, value := cell[0], cell[1], cell[2]

			board.SetCell(row, col, value)
		}
	}

	return hidedTarget
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

func (s *hideBacktracking) getCellShuffled(board *entities.Sudoku, r *rand.Rand) [][3]int {
	cellReference := make([][3]int, 0)

	for row := range board.Board.GetBoard() {
		for col := range board.Board.GetBoard()[row] {
			value := board.Board.GetCell(row, col)
			if value != 0 {
				cellReference = append(cellReference, [3]int{row, col, value})
			}
		}
	}

	r.Shuffle(len(cellReference), func(i, j int) {
		cellReference[i], cellReference[j] = cellReference[j], cellReference[i]
	})

	return cellReference
}
