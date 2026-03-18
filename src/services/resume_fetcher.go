package services

import (
	"context"

	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
)

type (
	resumeFetcher struct {
		sudokuRepository repository.SudokuRepository
	}
)

func NewResumeFetcher(sudokuRepository repository.SudokuRepository) domain.ResumeFetcher {
	return &resumeFetcher{sudokuRepository: sudokuRepository}
}

func (r *resumeFetcher) GetTotalSolvedByUser(ctx context.Context, userID vo.UUID) (map[entities.BoardSize]int, error) {
	return r.sudokuRepository.GetTotalSolvedByUser(ctx, userID)
}

func (r *resumeFetcher) GetTodaySolvedByUser(ctx context.Context, userID vo.UUID) ([]entities.GameResult, error) {
	solves, err := r.sudokuRepository.GetTodaySolvedByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return entities.ConvertSolvesToGameResults(solves), nil
}

func (r *resumeFetcher) GetBestTimesByUser(ctx context.Context, userID vo.UUID) ([]entities.GameResult, error) {
	solves, err := r.sudokuRepository.GetBestTimesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return entities.ConvertSolvesToGameResults(solves), nil
}
