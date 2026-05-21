package sudoku

import (
	"context"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/infrastructure/logging"
	"time"
)

type (
	RemoveUnfinishedAttemptsUseCase interface {
		Execute(ctx context.Context) error
	}

	removeUnfinishedAttemptsUseCase struct {
		sudokuRepository repository.SudokuRepository
	}
)

func NewRemoveUnfinishedAttemptsUseCase(
	sudokuRepository repository.SudokuRepository,
) RemoveUnfinishedAttemptsUseCase {
	return &removeUnfinishedAttemptsUseCase{
		sudokuRepository: sudokuRepository,
	}
}

func (r *removeUnfinishedAttemptsUseCase) Execute(ctx context.Context) error {
	today := time.Now().Truncate(24 * time.Hour)

	attempts, err := r.sudokuRepository.RemoveUnfinishedAttempts(ctx, today)
	if err != nil {
		return err
	}

	logging.Log(ctx).Info().
		Int64("attempts removed", attempts).
		Msg("attempts removed")

	return nil
}
