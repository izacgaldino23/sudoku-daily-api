package persistence

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"time"

	"github.com/uptrace/bun"
)

type (
	ISudokuRepository interface {
		WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
		Create(ctx context.Context, sudoku *entities.Sudoku) error
		GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error)
	}

	sudokuRepository struct {
		db bun.IDB
		tm repository.TransactionManager
	}
)

func NewSudokuRepository(db bun.IDB, transactionManager repository.TransactionManager) ISudokuRepository {
	return &sudokuRepository{
		db: db,
		tm: transactionManager,
	}
}

func (s *sudokuRepository) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.tm.WithinTransaction(ctx, fn)
}

func (s *sudokuRepository) GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error) {
	var sudokuResp Sudoku

	err := s.getExecutor(ctx).NewSelect().Model(&sudokuResp).Where("size = ? and date = ?", size, date).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return sudokuResp.ToDomain(), nil
}

func (s *sudokuRepository) Create(ctx context.Context, sudoku *entities.Sudoku) error {
	var sudokuModel = &Sudoku{}
	sudokuModel.FromDomain(sudoku)

	result, err := s.getExecutor(ctx).
		NewInsert().
		Model(sudokuModel).
		Column("id", "size", "difficulty", "board", "date").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}

func (s *sudokuRepository) getExecutor(ctx context.Context) bun.IDB {
	if tx, ok := extractTx(ctx); ok {
		return tx
	}

	return s.db
}
