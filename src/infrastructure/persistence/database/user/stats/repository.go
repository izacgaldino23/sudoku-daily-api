package stats

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/persistence/database/tx"

	"github.com/uptrace/bun"
)

type (
	userStatsRepository struct {
		txManager *tx.Manager
		db        *bun.DB
	}
)

func NewRepository(db *bun.DB) repository.UserStatsRepository {
	return &userStatsRepository{
		db:        db,
		txManager: tx.NewManager(db),
	}
}

func (u *userStatsRepository) GetByUserID(ctx context.Context, userID vo.UUID) (*entities.UserStats, error) {
	stats := Stats{}
	err := u.db.NewSelect().
		Model(&stats).
		Column("users.username").
		Join("JOIN users ON user_stats.user_id = users.id").
		Where("user_id = ?", userID.String()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return stats.ToDomain(), nil
}

func (u *userStatsRepository) GetOrCreate(ctx context.Context, userID vo.UUID) (*entities.UserStats, error) {
	stats, err := u.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if stats != nil {
		return stats, nil
	}

	stats = &entities.UserStats{
		ID:             vo.NewUUID(),
		UserID:         userID,
		CurrentStreak:  1,
		LongestStreak:  1,
		LastSolvedDate: time.Now(),
		TotalSolved:    1,
	}

	result, err := u.txManager.GetExecutor(ctx).
		NewInsert().
		Model(stats).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = result.RowsAffected()
	return stats, err
}

func (u *userStatsRepository) Update(ctx context.Context, stats *entities.UserStats) error {
	// update stats
	result, err := u.txManager.GetExecutor(ctx).
		NewUpdate().
		Model(stats).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}
