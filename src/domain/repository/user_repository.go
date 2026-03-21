package repository

import (
	"context"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}

type UserStatsRepository interface {
	GetByUserID(ctx context.Context, userID vo.UUID) (*entities.UserStats, error)
	GetOrCreate(ctx context.Context, userID vo.UUID) (*entities.UserStats, error)
	Update(ctx context.Context, stats *entities.UserStats) error
}
