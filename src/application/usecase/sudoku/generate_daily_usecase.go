package sudoku

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/logging"
)

type (
	GenerateDailyUseCase interface {
		Execute(ctx context.Context, size entities.BoardSize, date time.Time) (*entities.Sudoku, error)
	}

	sudokuGenerateDailyUseCase struct {
		txManager            repository.TransactionManager
		sudokuRepo           repository.SudokuRepository
		sudokuService        domain.SudokuGenerator
		sudokuFetcherService domain.SudokuDailyFetcher
	}
)

func NewSudokuGenerateDailyUseCase(
	txManager repository.TransactionManager,
	sudokuRepo repository.SudokuRepository,
	sudokuService domain.SudokuGenerator,
	sudokuFetchService domain.SudokuDailyFetcher,
) GenerateDailyUseCase {
	return &sudokuGenerateDailyUseCase{
		txManager:            txManager,
		sudokuRepo:           sudokuRepo,
		sudokuService:        sudokuService,
		sudokuFetcherService: sudokuFetchService,
	}
}

func (s *sudokuGenerateDailyUseCase) Execute(ctx context.Context, size entities.BoardSize, date time.Time) (*entities.Sudoku, error) {
	logging.Log(ctx).Info().Msgf("Generating sudoku for size %v", size)
	todayPuzzle, err := s.sudokuFetcherService.GetDaily(ctx, size)
	if err != nil && !errors.Is(err, pkg.ErrSudokuNotFound) {
		return nil, err
	}

	if todayPuzzle != nil && todayPuzzle.ID != "" {
		logging.Log(ctx).Info().Msgf("Sudoku for size %v already exists", size)
		return todayPuzzle, nil
	}

	sudoku, err := s.sudokuService.GenerateDaily(ctx, size, date)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sudoku for size %v: %w", size, err)
	}

	sudoku.Date = date
	sudoku.ID = vo.NewUUID()

	if err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		return s.sudokuRepo.Create(ctx, sudoku)
	}); err != nil {
		return nil, err
	}

	logging.Log(ctx).Info().Msgf("Sudoku for size %v generated", size)
	return sudoku, nil
}
