package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sudoku-daily-api/src/domain/entities"
	repository "sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"sudoku-daily-api/src/infrastructure/persistence/database/tx"

	"github.com/uptrace/bun"
)

type (
	sudokuRepository struct {
		txManager *tx.Manager
		db        *bun.DB
	}
)

func NewRepository(db *bun.DB) repository.SudokuRepository {
	return &sudokuRepository{
		db:        db,
		txManager: tx.NewManager(db),
	}
}

func (r *sudokuRepository) GetByDateAndSize(ctx context.Context, date time.Time, size entities.BoardSize) (*entities.Sudoku, error) {
	var sudokuResp Sudoku

	err := r.txManager.GetExecutor(ctx).NewSelect().Model(&sudokuResp).Where("size = ? and date = ?", size, date).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return sudokuResp.ToDomain(), nil
}

func (r *sudokuRepository) Create(ctx context.Context, sudoku *entities.Sudoku) error {
	var sudokuModel = &Sudoku{}
	sudokuModel.FromDomain(sudoku)

	result, err := r.txManager.GetExecutor(ctx).
		NewInsert().
		Model(sudokuModel).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (r *sudokuRepository) AddSolve(ctx context.Context, solve *entities.Solve) error {
	var solveModel = &Solve{}
	solveModel.FromDomain(solve)

	result, err := r.txManager.GetExecutor(ctx).
		NewInsert().
		Model(solveModel).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (r *sudokuRepository) GetTotalSolvedByUser(ctx context.Context, userID vo.UUID) (map[entities.BoardSize]int, error) {
	var results []sizeCount

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Table("solves").
		Column("sudokus.size").
		ColumnExpr("COUNT(*) AS total").
		Where("solves.user_id = ?", userID).
		Join("JOIN sudokus ON solves.sudoku_id = sudokus.id").
		Group("sudokus.size").
		Scan(ctx, &results)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	totalSolvesBySize := make(map[entities.BoardSize]int)
	for _, result := range results {
		totalSolvesBySize[entities.BoardSize(result.Size)] = result.Total
	}
	return totalSolvesBySize, nil
}

func (r *sudokuRepository) GetTodaySolvedByUser(ctx context.Context, userID vo.UUID) ([]entities.Solve, error) {
	var (
		today    = time.Now().Truncate(24 * time.Hour)
		tomorrow = today.Add(24 * time.Hour)
	)

	var solves = []Solve{}

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Model(&solves).
		Where("user_id = ? AND started_at >= ? AND started_at < ?", userID, today, tomorrow).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Solve, len(solves))
	for i, solve := range solves {
		result[i] = *solve.ToDomain()
	}

	return result, nil
}

func (r *sudokuRepository) GetBestTimesByUser(ctx context.Context, userID vo.UUID) ([]entities.Solve, error) {
	var solves []Solve

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Model(&solves).
		DistinctOn("sudokus.size").
		Join("JOIN sudokus ON solve.sudoku_id = sudokus.id").
		Where("solve.user_id = ? AND duration > 0", userID).
		Order("sudokus.size", "solve.duration ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Solve, len(solves))
	for i, solve := range solves {
		result[i] = *solve.ToDomain()
	}

	return result, nil
}

func (r *sudokuRepository) GetDailyLeaderboard(ctx context.Context, sudokuID vo.UUID, limit, offset int) ([]entities.Solve, bool, error) {
	var solves []Solve

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Column("users.username").
		Model(&solves).
		Join("JOIN users ON solve.user_id = users.id").
		Where("sudoku.id = ?", sudokuID).
		Order("solves.duration").
		Limit(limit + 1).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, false, err
	}

	hasNext := len(solves) > limit
	if hasNext {
		solves = solves[:limit]
	}

	result := make([]entities.Solve, len(solves)-1)
	for i, solve := range solves {
		result[i] = *solve.ToDomain()
	}

	return result, hasNext, nil
}

func (r *sudokuRepository) GetAllTimeBestLeaderboard(ctx context.Context, size entities.BoardSize, limit, offset int) ([]entities.Solve, bool, error) {
	var solves []Solve

	err := r.txManager.GetExecutor(ctx).
		NewSelect().
		Column("users.username").
		ExcludeColumn("duration").
		ColumnExpr("MIN(duration) AS duration").
		Model(&solves).
		Join("JOIN users ON solve.user_id = users.id").
		Where("size = ?", size).
		Group("user_id").
		Order("solves.duration").
		Limit(limit + 1).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, false, err
	}

	hasNext := len(solves) > limit
	if hasNext {
		solves = solves[:limit]
	}

	result := make([]entities.Solve, len(solves)-1)
	for i, solve := range solves {
		result[i] = *solve.ToDomain()
	}

	return result, hasNext, nil
}
