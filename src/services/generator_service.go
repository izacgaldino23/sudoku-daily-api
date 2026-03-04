package services

import (
	"fmt"
	"math/rand"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/helpers"
	"time"
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

func (s *sudokuGenerator) GenerateDaily(size entities.BoardSize, seed int64) (*entities.Sudoku, error) {
	sum := 0
	for i := 0; i < int(size); i++ {
		sum += i + 1
	}

	sudoku := entities.NewSudoku(size)

	r := rand.New(rand.NewSource(sudoku.Date.Unix()))
	start := time.Now()
	fmt.Printf("Start generating %v x %v sudoku at %v\n", size, size, start)

	filled := s.fillBacktracking.Fill(sudoku, r)
	if !filled {
		return nil, fmt.Errorf("failed to fill sudoku")
	}

	hide := s.hideBacktracking.Hide(sudoku, r)
	if !hide {
		return nil, fmt.Errorf("failed to hide sudoku")
	}

	fmt.Printf("Finish generating %v x %v sudoku at %v\n", size, size, time.Since(start).Seconds())

	return sudoku, nil
}
