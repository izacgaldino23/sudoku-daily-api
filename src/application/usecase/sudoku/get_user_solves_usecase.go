package sudoku

import (
	"context"
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	SudokuGetUserSolvesUseCase interface {
		Execute(ctx context.Context, userID vo.UUID) ([]entities.Solve, error)
	}

	getUserSolvesUseCase struct {
		sudokuRepository repository.SudokuRepository
	}
)

func NewSudokuGetUserSolvesUseCase(sudokuRepository repository.SudokuRepository) SudokuGetUserSolvesUseCase {
	return &getUserSolvesUseCase{sudokuRepository: sudokuRepository}
}

func (s *getUserSolvesUseCase) Execute(ctx context.Context, userID vo.UUID) ([]entities.Solve, error) {
	today := time.Now().Truncate(24 * time.Hour)

	solves, err := s.sudokuRepository.GetSolvesByUserAndDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	if len(solves) == 0 {
		return []entities.Solve{}, nil
	}

	return solves, nil
}
