package persistence

import (
	"context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"time"

	"github.com/uptrace/bun"
)

type (
	sudokuRepository struct {
		transactionManager
		db *bun.DB
	}
)

func NewSudokuRepository(db *bun.DB) repository.SudokuRepository {
	return &sudokuRepository{
		db:                 db,
		transactionManager: transactionManager{db: db},
	}
}

func (s *sudokuRepository) GetByDateAndSize(ctx context.Context, date time.Time, size int) (*entities.Sudoku, error) {
	var sudokuResp Sudoku

	err := s.GetExecutor(ctx).NewSelect().Model(&sudokuResp).Where("size = ? and date = ?", size, date).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return sudokuResp.ToDomain(), nil
}

func (s *sudokuRepository) Create(ctx context.Context, sudoku *entities.Sudoku) error {
	var sudokuModel = &Sudoku{}
	sudokuModel.FromDomain(sudoku)

	result, err := s.GetExecutor(ctx).
		NewInsert().
		Model(sudokuModel).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	return err
}
