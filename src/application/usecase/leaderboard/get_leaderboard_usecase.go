package leaderboard

import (
	"context"

	"sudoku-daily-api/src/domain/entities"
)

type (
	LeaderboardUsecase interface {
		Execute(ctx context.Context, params entities.LeaderboardSearchParams) (entities.Leaderboard, error)
	}
)
