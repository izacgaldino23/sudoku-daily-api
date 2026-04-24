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
	SudokuGenerateDailyUseCase interface {
		Execute(ctx context.Context) ([]entities.Sudoku, error)
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
) SudokuGenerateDailyUseCase {
	return &sudokuGenerateDailyUseCase{
		txManager:            txManager,
		sudokuRepo:           sudokuRepo,
		sudokuService:        sudokuService,
		sudokuFetcherService: sudokuFetchService,
	}
}

func (s *sudokuGenerateDailyUseCase) Execute(ctx context.Context) ([]entities.Sudoku, error) {
	var sudokuList []entities.Sudoku

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var puzzles []entities.Sudoku

	for boardSize := range entities.BoardSizes {
		logging.Log(ctx).Info().Msgf("Generating sudoku for size %v", boardSize)
		todayPuzzle, err := s.sudokuFetcherService.GetDaily(ctx, boardSize)
		if err != nil && !errors.Is(err, pkg.ErrSudokuNotFound) {
			return nil, err
		}

		if todayPuzzle != nil && todayPuzzle.ID != "" {
			sudokuList = append(sudokuList, *todayPuzzle)
			logging.Log(ctx).Info().Msgf("Sudoku for size %v already exists", boardSize)
			continue
		}

		sudoku, err := s.sudokuService.GenerateDaily(ctx, boardSize, today)
		if err != nil {
			return nil, fmt.Errorf("Failed to generate sudoku for size %v: %w", boardSize, err)
		}

		puzzles = append(puzzles, *sudoku)
		logging.Log(ctx).Info().Msgf("Sudoku for size %v generated", boardSize)
	}

	if err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, sudoku := range puzzles {
			sudoku.Date = today
			sudoku.ID = vo.NewUUID()

			err := s.sudokuRepo.Create(ctx, &sudoku)
			if err != nil {
				return err
			}

			sudokuList = append(sudokuList, sudoku)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return sudokuList, nil
}
