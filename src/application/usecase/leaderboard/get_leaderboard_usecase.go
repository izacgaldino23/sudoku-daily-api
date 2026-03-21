package leaderboard

import (
	"context"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
)

type (
	LeaderboardUsecase interface {
		Execute(ctx context.Context, params entities.LeaderboardSearchParams) (entities.Leaderboard, error)
	}

	leaderboardUsecase struct {
		UserStatsRepository repository.UserStatsRepository
	}
)
