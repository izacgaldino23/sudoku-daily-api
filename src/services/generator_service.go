package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/strategies"
	"sudoku-daily-api/src/infrastructure/logging"
)

type (
	sudokuGenerator struct {
		fillStrategy strategies.FillStrategy
		hideStrategy strategies.HideStrategy
	}
)

func NewGenerator(
	fillStrategy strategies.FillStrategy,
	hideStrategy strategies.HideStrategy,
) domain.SudokuGenerator {
	return &sudokuGenerator{
		fillStrategy: fillStrategy,
		hideStrategy: hideStrategy,
	}
}

func (s *sudokuGenerator) GenerateDaily(ctx context.Context, size entities.BoardSize, date time.Time) (*entities.Sudoku, error) {
	sudoku := entities.NewSudoku(size)

	r := rand.New(rand.NewSource(date.Unix()))

	filled := s.fillStrategy.Fill(sudoku, r)
	if !filled {
		return nil, fmt.Errorf("failed to fill sudoku")
	}
	logging.Log(ctx).Info().Msgf("Sudoku for size %v filled", size)

	cloned, err := sudoku.Board.Clone()
	if err != nil {
		return nil, err
	}
	sudoku.Solution = *cloned

	hide := s.hideStrategy.Hide(ctx, sudoku, r)
	if !hide {
		return nil, fmt.Errorf("failed to hide sudoku")
	}
	logging.Log(ctx).Info().Msgf("Sudoku for size %v hidden", size)

	return sudoku, nil
}
