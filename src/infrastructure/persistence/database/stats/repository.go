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
		Column("users.username", "stats.*").
		Join("JOIN users ON stats.user_id = users.id").
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
		return nil, r.txManager.HandleError(ctx, err)
	}

	if stats != nil {
		return stats, nil
	}

	newStats := &Stats{
		ID:             vo.NewUUID().String(),
		UserID:         userID.String(),
		CurrentStreak:  0,
		LongestStreak:  0,
		LastSolvedDate: time.Now(),
		TotalSolved:    0,
	}

	result, err := r.txManager.GetExecutor(ctx).
		NewInsert().
		Model(newStats).
		Exec(ctx)
	if err != nil {
		return nil, r.txManager.HandleError(ctx, err)
	}

	_, err = result.RowsAffected()
	return newStats.ToDomain(), r.txManager.HandleError(ctx, err)
}

func (r *userStatsRepository) Update(ctx context.Context, stats *entities.UserStats) error {
	statsModel := &Stats{}
	statsModel.FromDomain(stats)

	result, err := r.txManager.GetExecutor(ctx).
		NewUpdate().
		Model(statsModel).
		Where("id = ?", stats.ID).
		Exec(ctx)
	if err != nil {
		return r.txManager.HandleError(ctx, err)
	}

	_, err = result.RowsAffected()
	return r.txManager.HandleError(ctx, err)
}

func (r *userStatsRepository) GetBestStreakLeaderboard(ctx context.Context, limit int, offset int, filterDate time.Time) ([]entities.UserStats, bool, error) {
	var stats []Stats

	currentDate := filterDate.Truncate(24 * time.Hour)
	dayBefore := currentDate.AddDate(0, 0, -1)

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Column("users.username", "stats.*").
		Model(&stats).
		Join("JOIN users ON stats.user_id = users.id").
		Where("last_solved_date >= ? AND last_solved_date <= ?", dayBefore, currentDate).
		Order("stats.longest_streak DESC").
		Limit(limit + 1).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, false, r.txManager.HandleError(ctx, err)
	}

	hasNext := len(stats) > limit
	if hasNext {
		stats = stats[:limit]
	}

	if len(stats) == 0 {
		return nil, false, nil
	}

	result := make([]entities.UserStats, len(stats))
	for i, stat := range stats {
		result[i] = *stat.ToDomain()
	}

	return result, hasNext, nil
}

func (r *userStatsRepository) GetTotalSolvesLeaderboard(ctx context.Context, limit int, offset int) ([]entities.UserStats, bool, error) {
	var stats []Stats

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Model(&stats).
		Column("users.username", "stats.*").
		Join("JOIN users ON stats.user_id = users.id").
		Order("stats.total_solved DESC").
		Limit(limit + 1).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, false, r.txManager.HandleError(ctx, err)
	}

	hasNext := len(stats) > limit
	if hasNext {
		stats = stats[:limit]
	}

	if len(stats) == 0 {
		return nil, false, nil
	}

	result := make([]entities.UserStats, len(stats))
	for i, stat := range stats {
		result[i] = *stat.ToDomain()
	}

	return result, hasNext, nil
}

func (r *userStatsRepository) ResetStrikes(ctx context.Context, today time.Time) (int64, error) {
	dayBeforeYesterday := today.Truncate(24*time.Hour).AddDate(0, 0, -2)

	result, err := r.txManager.GetExecutor(ctx).
		NewUpdate().
		Model(&Stats{}).
		Where("last_solved_date <= ?", dayBeforeYesterday).
		Exec(ctx)
	if err != nil {
		return 0, r.txManager.HandleError(ctx, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, r.txManager.HandleError(ctx, err)
	}

	return count, nil
}
