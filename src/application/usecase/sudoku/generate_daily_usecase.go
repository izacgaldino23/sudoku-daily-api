package sudoku

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
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

	if err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		for boardSize := range entities.BoardSizes {
			todayPuzzle, err := s.sudokuFetcherService.GetDaily(ctx, boardSize)
			if err != nil && !errors.Is(err, pkg.ErrSudokuNotFound) {
				return err
			}

			if todayPuzzle != nil && todayPuzzle.ID != "" {
				sudokuList = append(sudokuList, *todayPuzzle)
				log.Warn().Msgf("Sudoku for size %v already exists", boardSize)
				continue
			}

			sudoku, err := s.sudokuService.GenerateDaily(boardSize, today)
			if err != nil {
				return fmt.Errorf("Failed to generate sudoku for size %v: %w", boardSize, err)
			}

			sudoku.Date = today
			sudoku.ID = vo.NewUUID()

			err = s.sudokuRepo.Create(ctx, sudoku)
			if err != nil {
				return err
			}

			sudokuList = append(sudokuList, *sudoku)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return sudokuList, nil
}
