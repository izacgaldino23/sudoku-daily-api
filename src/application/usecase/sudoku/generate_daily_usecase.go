package sudoku

import (
	"context"
	"fmt"
	"time"

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
		txManager     repository.TransactionManager
		sudokuRepo    repository.SudokuRepository
		sudokuService domain.SudokuGenerator
	}
)

func NewSudokuGenerateDailyUseCase(
	txManager repository.TransactionManager,
	sudokuRepo repository.SudokuRepository,
	sudokuService domain.SudokuGenerator,
) SudokuGenerateDailyUseCase {
	return &sudokuGenerateDailyUseCase{
		txManager:     txManager,
		sudokuRepo:    sudokuRepo,
		sudokuService: sudokuService,
	}
}

func (s *sudokuGenerateDailyUseCase) Execute(ctx context.Context) ([]entities.Sudoku, error) {
	var sudokuList []entities.Sudoku

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if err := s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		for boardSize := range entities.BoardSizes {
			sudoku, err := s.sudokuService.GenerateDaily(boardSize, today.UnixNano())
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
