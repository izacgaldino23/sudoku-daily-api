package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size int, sessionID vo.UUID) (sudoku *entities.Sudoku, token string, err error)
	}

	sudokuGetDailyUseCase struct {
		tokenService  domain.TokenService
		sudokuFetcher domain.SudokuDailyFetcher
	}
)

func NewSudokuGetDailyUseCase(
	tokenService domain.TokenService,
	sudokuFetcher domain.SudokuDailyFetcher,
) ISudokuGetDailyUseCase {
	return &sudokuGetDailyUseCase{
		tokenService:  tokenService,
		sudokuFetcher: sudokuFetcher,
	}
}

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size int, sessionID vo.UUID) (*entities.Sudoku, string, error) {
	_, ok := entities.BoardSizes[entities.BoardSize(size)]
	if !ok {
		log.Ctx(ctx).Error().Msgf("Invalid size: %d", size)
		return nil, "", pkg.ErrQueryParamInvalid
	}

	sudoku, err := s.sudokuFetcher.GetDaily(ctx, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", pkg.ErrNotFound
		}
		return nil, "", err
	}

	token, err := s.generateToken(sessionID, sudoku)
	if err != nil {
		return nil, "", err
	}

	return sudoku, token, nil
}

func (s *sudokuGetDailyUseCase) generateToken(sessionID vo.UUID, sudoku *entities.Sudoku) (string, error) {
	tomorrow := sudoku.Date.AddDate(0, 0, 1)

	sessionToken := &entities.SessionToken{
		Date:       sudoku.Date.Format(time.DateOnly),
		Size:       int(sudoku.Size),
		SessionID:  sessionID,
		StartedAt:  time.Now(),
		ExpiresAt:  tomorrow,
		
	}

	token, err := s.tokenService.GenerateJWTToken(sessionToken.ToMap())
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate session token")
		return "", err
	}

	return token, nil
}
