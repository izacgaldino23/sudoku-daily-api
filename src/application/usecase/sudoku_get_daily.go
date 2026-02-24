package usecase

import (
	"context"
	"sudoku-daily-api/src/adapters/repository"
	"sudoku-daily-api/src/domain/entities"
	"time"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size int) (*entities.Sudoku, error)
	}

	sudokuGetDailyUseCase struct {
		repository repository.ISudokuRepository
	}
)

func NewSudokuGetDailyUseCase(repository repository.ISudokuRepository) ISudokuGetDailyUseCase {
	return &sudokuGetDailyUseCase{repository: repository}
}

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size int) (*entities.Sudoku, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// Add in cache

	return s.repository.GetByDateAndSize(ctx, today, size)
}