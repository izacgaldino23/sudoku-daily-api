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

func (r *sudokuRepository) GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error) {
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
	query := `SELECT sudokus.size, COUNT(*) AS total 
		FROM solves 
		JOIN sudokus ON solves.sudoku_id = sudokus.id 
		WHERE solves.user_id = ? 
		GROUP BY sudokus.size`
	err := r.txManager.GetExecutor(ctx).NewRaw(query, userID).Scan(ctx, &results)
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
	query := `SELECT DISTINCT ON (sudokus.size) solves.* 
		FROM solves 
		JOIN sudokus ON solves.sudoku_id = sudokus.id 
		WHERE solves.user_id = ? AND solves.duration > 0 
		ORDER BY sudokus.size, solves.duration ASC`
	err := r.txManager.GetExecutor(ctx).NewRaw(query, userID).Scan(ctx, &solves)
	if err != nil {
		return nil, err
	}

	result := make([]entities.Solve, len(solves))
	for i, solve := range solves {
		result[i] = *solve.ToDomain()
	}

	return result, nil
}
