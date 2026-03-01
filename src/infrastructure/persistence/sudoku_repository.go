package persistence

import (
	"context"
	"database/sql"
	"sudoku-daily-api/src/domain/entities"
	"time"

	"github.com/uptrace/bun"
)

type (
	ISudokuRepository interface {
		GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error)
		Create(ctx context.Context, sudoku *entities.Sudoku) error
		RunTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
	}

	sudokuRepository struct {
		db bun.IDB
	}
)

func NewSudokuRepository(db bun.IDB) ISudokuRepository {
	return &sudokuRepository{db: db}
}

func (s *sudokuRepository) GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error) {
	var sudokuResp Sudoku

	err := s.db.NewSelect().Model(&sudokuResp).Where("size = ? and date = ?", size, date).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return sudokuResp.ToDomain(), nil
}

func (s *sudokuRepository) RunTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return s.db.RunInTx(ctx, &sql.TxOptions{}, fn)
}

func (s *sudokuRepository) Create(ctx context.Context, sudoku *entities.Sudoku) error {
	var sudokuModel Sudoku
	sudokuModel.FromDomain(sudoku)

	result, err := s.db.NewInsert().Model(sudokuModel).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}
