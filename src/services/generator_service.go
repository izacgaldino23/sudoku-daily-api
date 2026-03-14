package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/strategies"
	"time"
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

func (s *sudokuGenerator) GenerateDaily(size entities.BoardSize, seed int64) (*entities.Sudoku, error) {
	sum := 0
	for i := 0; i < int(size); i++ {
		sum += i + 1
	}

	sudoku := entities.NewSudoku(size)

	r := rand.New(rand.NewSource(sudoku.Date.Unix()))
	start := time.Now()
	fmt.Printf("Start generating %v x %v sudoku at %v\n", size, size, start)

	filled := s.fillStrategy.Fill(sudoku, r)
	if !filled {
		return nil, fmt.Errorf("failed to fill sudoku")
	}

	if err := deepCopy(&sudoku.Board, &sudoku.Solution); err != nil {
		return nil, err
	}

	hide := s.hideStrategy.Hide(sudoku, r)
	if !hide {
		return nil, fmt.Errorf("failed to hide sudoku")
	}

	fmt.Printf("Finish generating %v x %v sudoku in %vms\n", size, size, time.Since(start).Milliseconds())

	return sudoku, nil
}

func deepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return err
	}
	return nil
}
