package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	sudokuDailyFetcher struct {
		cache            domain.Cache
		sudokuRepository repository.SudokuRepository
	}
)

func NewSudokuDailyFetcher(cache domain.Cache, sudokuRepository repository.SudokuRepository) domain.SudokuDailyFetcher {
	return &sudokuDailyFetcher{
		cache:            cache,
		sudokuRepository: sudokuRepository,
	}
}

func (s *sudokuDailyFetcher) GetDaily(ctx context.Context, size entities.BoardSize) (*entities.Sudoku, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	cacheKey := fmt.Sprintf("sudoku-%d", size)
	if value, ok := s.cache.Get(cacheKey); ok {
		sudoku := value.(*entities.Sudoku)

		if isSameDate(sudoku.Date, today) {
			return sudoku, nil
		}
	}

	sudoku, err := s.sudokuRepository.GetByDateAndSize(ctx, today, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSudokuNotFound
		}
		return nil, err
	}

	s.cache.Set(cacheKey, sudoku)

	return sudoku, nil
}

func (s *sudokuDailyFetcher) GetByDateAndSize(ctx context.Context, date time.Time, size entities.BoardSize) (*entities.Sudoku, error) {
	return s.sudokuRepository.GetByDateAndSize(ctx, date, size)
}

func (s *sudokuDailyFetcher) GetSolveByIDAndUser(ctx context.Context, sudokuID, userID vo.UUID) (*entities.Solve, error) {
	solve, err := s.sudokuRepository.GetSolveByIDAndUser(ctx, sudokuID, userID)
	if err != nil {
		if errors.Is(err, pkg.ErrSolutionNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return solve, nil
}

func isSameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
