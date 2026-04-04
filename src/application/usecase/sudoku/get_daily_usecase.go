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
	"sudoku-daily-api/src/infrastructure/logging"

	"github.com/rs/zerolog/log"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size entities.BoardSize, userID vo.UUID) (sudoku *entities.Sudoku, playToken string, err error)
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

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size entities.BoardSize, userID vo.UUID) (*entities.Sudoku, string, error) {
	boardSize := entities.BoardSize(size)

	sudoku, err := s.sudokuFetcher.GetDaily(ctx, boardSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", pkg.ErrNotFound
		}
		return nil, "", err
	}

	if !userID.IsEmpty() {
		// validate if user has already played the game
		if _, err = s.sudokuFetcher.GetSolveByIDAndUser(ctx, sudoku.ID, userID); err != nil {
			logging.Log(ctx).Info().Err(err).Msg("user has already played the game")
			if errors.Is(err, pkg.ErrNotFound) {
				return nil, "", pkg.ErrAlreadyPlayed
			}
		}
	}

	sessionID := app_context.GetSessionIDFromContext(ctx)

	token, err := s.generateToken(sessionID, sudoku)
	if err != nil {
		return nil, "", err
	}

	return sudoku, token, nil
}

func (s *sudokuGetDailyUseCase) generateToken(sessionID vo.UUID, sudoku *entities.Sudoku) (string, error) {
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
