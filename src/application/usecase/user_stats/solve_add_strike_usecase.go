package user_stats_usecase

import (
	"context"
	"time"

	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	SolveAddStrikeUseCase interface {
		Execute(ctx context.Context, userID vo.UUID, solveDate time.Time) error
	}

	solveAddStrikeUseCase struct {
		userStatsRepository repository.UserStatsRepository
	}
)

func NewSolveAddStrikeUseCase(userStatsRepository repository.UserStatsRepository) SolveAddStrikeUseCase {
	return &solveAddStrikeUseCase{
		userStatsRepository: userStatsRepository,
	}
}

func (s *solveAddStrikeUseCase) Execute(ctx context.Context, userID vo.UUID, solveDate time.Time) error {
	solveDateOnly := solveDate.Truncate(24 * time.Hour)

	stats, err := s.userStatsRepository.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	if !stats.LastSolvedDate.IsZero() {
		yesterday := solveDateOnly.AddDate(0, 0, -1)
		if stats.LastSolvedDate.Equal(solveDateOnly) {
			return nil
		} else if stats.LastSolvedDate.Equal(yesterday) {
			stats.CurrentStreak++
		} else {
			stats.CurrentStreak = 1
		}
	}

	if stats.CurrentStreak > stats.LongestStreak {
		stats.LongestStreak = stats.CurrentStreak
	}
	
	stats.LastSolvedDate = solveDateOnly
	stats.TotalSolved++

	if err = s.userStatsRepository.Update(ctx, stats); err != nil {
		return err
	}

	return nil
}
