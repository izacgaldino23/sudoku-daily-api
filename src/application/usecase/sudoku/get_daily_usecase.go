package sudoku

import (
	"context"
	"database/sql"
	"errors"
	"sudoku-daily-api/src/domain/repository"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"

	"github.com/rs/zerolog/log"
)

type (
	ISudokuGetDailyUseCase interface {
		Execute(ctx context.Context, size entities.BoardSize, userID vo.UUID) (sudoku *entities.Sudoku, playToken string, startedAt time.Time, err error)
	}

	sudokuGetDailyUseCase struct {
		tokenService     domain.TokenService
		sudokuFetcher    domain.SudokuDailyFetcher
		sudokuRepository repository.SudokuRepository
	}
)

func NewSudokuGetDailyUseCase(
	tokenService domain.TokenService,
	sudokuFetcher domain.SudokuDailyFetcher,
	sudokuRepository repository.SudokuRepository,
) ISudokuGetDailyUseCase {
	return &sudokuGetDailyUseCase{
		tokenService:     tokenService,
		sudokuFetcher:    sudokuFetcher,
		sudokuRepository: sudokuRepository,
	}
}

func (s *sudokuGetDailyUseCase) Execute(ctx context.Context, size entities.BoardSize, userID vo.UUID) (*entities.Sudoku, string, time.Time, error) {
	var (
		boardSize = entities.BoardSize(size)
		startedAt = time.Now()
	)

	sudoku, err := s.sudokuFetcher.GetDaily(ctx, boardSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", startedAt, pkg.ErrSudokuNotFound
		}
		return nil, "", startedAt, err
	}

	// get attempt for current solve
	solve, err := s.sudokuFetcher.GetSolveByIDAndUser(ctx, sudoku.ID, userID)
	if err != nil {
		if !errors.Is(err, pkg.ErrSolutionNotFound) {
			return nil, "", startedAt, err
		}
	}

	if solve != nil && !solve.ID.IsEmpty() {
		// Is solved
		if solve.Duration != 0 {
			return nil, "", startedAt, pkg.ErrAlreadyPlayed
		}

		startedAt = solve.StartedAt
	}

	var token string
	token, err = s.generateTokenWithUser(userID, sudoku, startedAt)
	if err != nil {
		return nil, "", startedAt, err
	}

	if err = s.sudokuRepository.AddAttempt(ctx, &entities.Solve{
		ID:        vo.NewUUID(),
		SudokuID:  sudoku.ID,
		UserID:    userID,
		StartedAt: startedAt,
		CreatedAt: time.Now(),
		Size:      int(size),
	}); err != nil {
		return nil, "", startedAt, err
	}

	return sudoku, token, startedAt, nil
}

func (s *sudokuGetDailyUseCase) generateTokenWithUser(userID vo.UUID, sudoku *entities.Sudoku, startedAt time.Time) (string, error) {
	tomorrow := sudoku.Date.AddDate(0, 0, 1)

	playToken := &entities.PlayToken{
		Date:      sudoku.Date.Format(time.DateOnly),
		Size:      sudoku.Size,
		UserID:    userID,
		SudokuID:  sudoku.ID,
		StartedAt: startedAt,
		ExpiresAt: tomorrow,
	}

	secondsUntilTomorrow := int(time.Until(tomorrow).Seconds())
	if secondsUntilTomorrow < 0 {
		secondsUntilTomorrow = 0
	}

	token, err := s.tokenService.GenerateJWTToken(playToken.ToMap(), &secondsUntilTomorrow)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate user token")
		return "", err
	}

	return token, nil
}
