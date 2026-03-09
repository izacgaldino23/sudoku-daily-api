package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"time"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size int) (*entities.Sudoku, error)
	}

	sudokuGetDailyUseCase struct {
		repository repository.SudokuRepository
		cache      domain.Cache
	}
)

func NewSudokuGetDailyUseCase(
	repository repository.SudokuRepository,
	cache domain.Cache,
) ISudokuGetDailyUseCase {
	return &sudokuGetDailyUseCase{
		repository: repository,
		cache:      cache,
	}
}

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size int) (*entities.Sudoku, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	cacheKey := fmt.Sprintf("sudoku-%d", size)
	if value, ok := s.cache.Get(cacheKey); ok {
		sudoku := value.(*entities.Sudoku)

		if isSameDate(sudoku.Date, today) {
			return sudoku, nil
		}
	}

	board, err := s.repository.GetByDateAndSize(ctx, today, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}
		return nil, err
	}

	s.cache.Set(cacheKey, board)

	return board, nil
}

func isSameDate(a, b time.Time) bool { 
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day() 
}