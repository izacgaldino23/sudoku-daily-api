package services

import (
	"math/rand"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
)

type (
	sudokuGenerator struct {
		fillBacktracking helpers.FillBacktracking
		hideBacktracking helpers.HideBacktracking
	}
)

func NewGenerator(
	fillBacktracking helpers.FillBacktracking,
	hideBacktracking helpers.HideBacktracking,
) domain.SudokuGenerator {
	return &sudokuGenerator{
		fillBacktracking: fillBacktracking,
		hideBacktracking: hideBacktracking,
	}
}

func (s *sudokuGenerator) GenerateDaily(size int, seed int64) *entities.Sudoku {
	sum := 0
	for i := 0; i < size; i++ {
		sum += i + 1
	}

	board := &entities.Sudoku{
		Size:  size,
		Board: make([][]int, size),
	}

	// initialize board
	for i := range board.Board {
		board.Board[i] = make([]int, size)
	}

	r := rand.New(rand.NewSource(board.Date.Unix()))

	s.fillBacktracking.Fill(board, r)
	s.hideBacktracking.Hide(board, r)

	return board
}
