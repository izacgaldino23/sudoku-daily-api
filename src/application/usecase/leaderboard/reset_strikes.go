package leaderboard_usecase

import (
	"context"
	"time"

	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/infrastructure/logging"
)

type (
	ResetStrikesUseCase interface {
		Execute(ctx context.Context, date time.Time) error
	}

	resetStrikesUseCase struct {
		userStatsRepository repository.UserStatsRepository
	}
)

func NewResetStrikesUseCase(userStatsRepository repository.UserStatsRepository) ResetStrikesUseCase {
	return &resetStrikesUseCase{userStatsRepository: userStatsRepository}
}

func (r *resetStrikesUseCase) Execute(ctx context.Context, date time.Time) error {
	reset, err := r.userStatsRepository.ResetStrikes(ctx, date)
	if err != nil {
		return err
	}

	logging.Log(ctx).Info().
		Int64("reset", reset).
		Msg("reset")

	return nil
}
