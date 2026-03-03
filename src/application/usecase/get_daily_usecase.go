package usecase

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/infrastructure/persistence"
	"time"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size int) (*entities.Sudoku, error)
	}

	sudokuGetDailyUseCase struct {
		repository persistence.ISudokuRepository
	}
)

func NewSudokuGetDailyUseCase(repository persistence.ISudokuRepository) ISudokuGetDailyUseCase {
	return &sudokuGetDailyUseCase{repository: repository}
}

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size int) (*entities.Sudoku, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// TODO: Add in cache
	board, err := s.repository.GetByDateAndSize(ctx, today, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrorNotFound
		}
		return nil, err
	}

	return board, nil
}
