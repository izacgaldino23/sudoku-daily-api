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

func (r *userStatsRepository) GetByUserID(ctx context.Context, userID vo.UUID) (*entities.UserStats, error) {
	stats := Stats{}
	err := r.txManager.GetExecutor(ctx).
		NewSelect().
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

func (r *userStatsRepository) GetOrCreate(ctx context.Context, userID vo.UUID) (*entities.UserStats, error) {
	stats, err := r.GetByUserID(ctx, userID)
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

	result, err := r.txManager.GetExecutor(ctx).
		NewInsert().
		Model(stats).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = result.RowsAffected()
	return stats, err
}

func (r *userStatsRepository) Update(ctx context.Context, stats *entities.UserStats) error {
	// update stats
	result, err := r.txManager.GetExecutor(ctx).
		NewUpdate().
		Model(stats).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (r *userStatsRepository) GetBestStreakLeaderboard(ctx context.Context, limit int, offset int, filterDate time.Time) ([]entities.UserStats, bool, error) {
	var stats []Stats

	currentDate := filterDate.Truncate(24 * time.Hour)
	dayBefore := currentDate.AddDate(0, 0, -1)

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Column("users.username").
		Model(&stats).
		Join("JOIN users ON solve.user_id = users.id").
		Where("last_solved_date >= ? AND last_solved_date <= ?", dayBefore, currentDate).
		Order("user_stats.longest_streak DESC").
		Limit(limit + 1).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, false, err
	}

	hasNext := len(stats) > limit
	if hasNext {
		stats = stats[:limit]
	}

	result := make([]entities.UserStats, len(stats)-1)
	for i, stat := range stats {
		result[i] = *stat.ToDomain()
	}

	return result, hasNext, nil
}

func (r *userStatsRepository) GetTotalSolvesLeaderboard(ctx context.Context, limit int, offset int) ([]entities.UserStats, bool, error) {
	var stats []Stats

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Column("users.username").
		Model(&stats).
		Join("JOIN users ON user_stats.user_id = users.id").
		Order("user_stats.total_solved DESC").
		Limit(limit + 1).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, false, err
	}

	hasNext := len(stats) > limit
	if hasNext {
		stats = stats[:limit]
	}

	result := make([]entities.UserStats, len(stats)-1)
	for i, stat := range stats {
		result[i] = *stat.ToDomain()
	}

	return result, hasNext, nil
}
