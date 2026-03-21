package user_stats_usecase

import (
	"context"
	"time"

	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	SolveAddStrikeUseCase interface {
		Execute(ctx context.Context, userID vo.UUID) error
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

func (s *solveAddStrikeUseCase) Execute(ctx context.Context, userID vo.UUID) error {
	today := time.Now().Truncate(time.Hour * 24)

	// get current stats
	stats, err := s.userStatsRepository.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	if stats.LastSolvedDate.Equal(today) {
		stats.CurrentStreak++
	} else {
		if stats.CurrentStreak > stats.LongestStreak {
			stats.LongestStreak = stats.CurrentStreak
		}
		stats.CurrentStreak = 1
	}
	
	stats.LastSolvedDate = time.Now()
	stats.TotalSolved++

	// update stats
	if err = s.userStatsRepository.Update(ctx, stats); err != nil {
		return err
	}

	return nil
}
