package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size int, sessionID vo.UUID) (sudoku *entities.Sudoku, token string, err error)
	}

	sudokuGetDailyUseCase struct {
		repository   repository.SudokuRepository
		tokenService domain.TokenService
		cache        domain.Cache
	}
)

func NewSudokuGetDailyUseCase(
	repository repository.SudokuRepository,
	tokenService domain.TokenService,
	cache domain.Cache,
) ISudokuGetDailyUseCase {
	return &sudokuGetDailyUseCase{
		repository:   repository,
		tokenService: tokenService,
		cache:        cache,
	}
}

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size int, sessionID vo.UUID) (*entities.Sudoku, string, error) {
	_, ok := entities.BoardSizes[entities.BoardSize(size)]
	if !ok {
		log.Ctx(ctx).Error().Msgf("Invalid size: %d", size)
		return nil, "", pkg.ErrQueryParamInvalid
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	cacheKey := fmt.Sprintf("sudoku-%d", size)
	if value, ok := s.cache.Get(cacheKey); ok {
		sudoku := value.(*entities.Sudoku)

		if isSameDate(sudoku.Date, today) {
			return sudoku, "", nil
		}
	}

	sudoku, err := s.repository.GetByDateAndSize(ctx, today, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", pkg.ErrNotFound
		}
		return nil, "", err
	}

	s.cache.Set(cacheKey, sudoku)

	token, err := s.generateToken(sessionID, sudoku)
	if err != nil {
		return nil, "", err
	}

	return sudoku, token, nil
}

func (s *sudokuGetDailyUseCase) generateToken(sessionID vo.UUID, sudoku *entities.Sudoku) (string, error) {
	tomorrow := sudoku.Date.AddDate(0, 0, 1)

	token, err := s.tokenService.GenerateJWTToken(map[string]any{
		"session_id": sessionID,
		"date":       sudoku.Date.Format(time.DateOnly),
		"start_time": time.Now().Unix(),
		"exp":        tomorrow.Unix(),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate session token")
		return "", err
	}

	return token, nil
}

func isSameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}
