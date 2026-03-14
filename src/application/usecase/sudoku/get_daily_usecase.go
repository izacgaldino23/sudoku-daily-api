package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size int) (sudoku *entities.Sudoku, playToken string, sessionID vo.UUID, err error)
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

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size int) (*entities.Sudoku, string, vo.UUID, error) {
	_, ok := entities.BoardSizes[entities.BoardSize(size)]
	if !ok {
		log.Ctx(ctx).Error().Msgf("Invalid size: %d", size)
		return nil, "", "", pkg.ErrQueryParamInvalid
	}

	sudoku, err := s.sudokuFetcher.GetDaily(ctx, size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", "", pkg.ErrNotFound
		}
		return nil, "", "", err
	}

	sessionID := app_context.GetSessionIDFromContext(ctx)

	token, sessionID, err := s.generateToken(sessionID, sudoku)
	if err != nil {
		return nil, "", "", err
	}

	return sudoku, token, sessionID, nil
}

func (s *sudokuGetDailyUseCase) generateToken(sessionID vo.UUID, sudoku *entities.Sudoku) (string, vo.UUID, error) {
	tomorrow := sudoku.Date.AddDate(0, 0, 1)

	if sessionID == "" {
		sessionID = vo.NewUUID()
	}

	playToken := &entities.PlayToken{
		Date:      sudoku.Date.Format(time.DateOnly),
		Size:      int(sudoku.Size),
		SessionID: sessionID,
		StartedAt: time.Now(),
		ExpiresAt: tomorrow,
	}

	// seconds until tomorrow
	
	secondsUntilTomorrow := int(time.Until(tomorrow).Seconds())
	if secondsUntilTomorrow < 0 {
		secondsUntilTomorrow = 0
	}

	token, err := s.tokenService.GenerateJWTToken(playToken.ToMap(), &secondsUntilTomorrow)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate session token")
		return "", "", err
	}

	return token, sessionID, nil
}
