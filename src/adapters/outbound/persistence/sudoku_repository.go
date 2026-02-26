package persistence

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
	"time"

	"github.com/uptrace/bun"
)

type (
	ISudokuRepository interface {
		GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error)
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