package usecase

import (
	"context"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/persistence"
	"time"
)

type (
	ISudokuGenerateAllUseCase interface {
		Execute(ctx context.Context) ([]entities.Sudoku, error)
	}

	sudokuGenerateAllUseCase struct {
		repository    persistence.ISudokuRepository
		sudokuService domain.SudokuGenerator
	}
)

func NewSudokuGenerateAllUseCase(
	repository persistence.ISudokuRepository,
	sudokuService domain.SudokuGenerator,
) ISudokuGenerateAllUseCase {
	return &sudokuGenerateAllUseCase{
		repository:    repository,
		sudokuService: sudokuService,
	}
}

func (s *sudokuGenerateAllUseCase) Execute(ctx context.Context) ([]entities.Sudoku, error) {
	var sudokuList []entities.Sudoku

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if err := s.repository.WithinTransaction(ctx, func(ctx context.Context) error {
		for boardSize := range entities.BoardSizes {
			sudoku := s.sudokuService.GenerateDaily(int(boardSize), today.UnixNano())

			sudoku.Date = today
			sudoku.ID = string(vo.NewUUID())

			err := s.repository.Create(ctx, sudoku)
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
