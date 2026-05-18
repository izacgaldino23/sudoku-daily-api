package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"

	"github.com/rs/zerolog/log"
)

type (
	ISudokuGetDailyForGuestUseCase interface {
		Execute(ctx context.Context, size entities.BoardSize) (sudoku *entities.Sudoku, playToken string, sessionID vo.UUID, err error)
	}

	sudokuGetDailyFroGuestUseCase struct {
		tokenService  domain.TokenService
		sudokuFetcher domain.SudokuDailyFetcher
	}
)

func NewSudokuGetDailyForGuestUseCase(
	tokenService domain.TokenService,
	sudokuFetcher domain.SudokuDailyFetcher,
) ISudokuGetDailyForGuestUseCase {
	return &sudokuGetDailyFroGuestUseCase{
		tokenService:  tokenService,
		sudokuFetcher: sudokuFetcher,
	}
}

func (s *sudokuGetDailyFroGuestUseCase) Execute(ctx context.Context, size entities.BoardSize) (*entities.Sudoku, string, vo.UUID, error) {
	boardSize := entities.BoardSize(size)

	sudoku, err := s.sudokuFetcher.GetDaily(ctx, boardSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", "", pkg.ErrSudokuNotFound
		}
		return nil, "", "", err
	}

	sessionID := app_context.GetSessionIDFromContext(ctx)

	var playToken string
	playToken, err = s.generateToken(sessionID, sudoku)
	if err != nil {
		return nil, "", "", err
	}

	return sudoku, playToken, sessionID, nil
}

func (s *sudokuGetDailyFroGuestUseCase) generateToken(sessionID vo.UUID, sudoku *entities.Sudoku) (string, error) {
	tomorrow := sudoku.Date.AddDate(0, 0, 1)

	playToken := &entities.PlayToken{
		Date:      sudoku.Date.Format(time.DateOnly),
		Size:      sudoku.Size,
		SessionID: sessionID,
		SudokuID:  sudoku.ID,
		StartedAt: time.Now(),
		ExpiresAt: tomorrow,
	}

	secondsUntilTomorrow := int(time.Until(tomorrow).Seconds())
	if secondsUntilTomorrow < 0 {
		secondsUntilTomorrow = 0
	}

	token, err := s.tokenService.GenerateJWTToken(playToken.ToMap(), &secondsUntilTomorrow)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate session token")
		return "", err
	}

	return token, nil
}
