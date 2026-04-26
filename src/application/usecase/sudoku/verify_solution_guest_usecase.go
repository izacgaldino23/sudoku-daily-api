package sudoku

import (
	"context"
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"

	"github.com/rs/zerolog/log"
)

type (
	SudokuVerifySolutionGuestUseCase interface {
		Execute(ctx context.Context, sudoku *entities.Solve, playToken string, finished time.Time) (bool, error)
	}

	sudokuVerifySolutionGuestUseCase struct {
		sudokuRepo    repository.SudokuRepository
		tokenService  domain.TokenService
		sudokuFetcher domain.SudokuDailyFetcher
		txManager     repository.TransactionManager
	}
)

func NewSudokuVerifySolutionGuestUseCase(
	sudokuRepo repository.SudokuRepository,
	tokenService domain.TokenService,
	sudokuFetcher domain.SudokuDailyFetcher,
	txManager repository.TransactionManager,
) SudokuVerifySolutionGuestUseCase {
	return &sudokuVerifySolutionGuestUseCase{
		sudokuRepo:    sudokuRepo,
		tokenService:  tokenService,
		sudokuFetcher: sudokuFetcher,
		txManager:     txManager,
	}
}

func (s *sudokuVerifySolutionGuestUseCase) Execute(ctx context.Context, input *entities.Solve, token string, finished time.Time) (bool, error) {
	claims, err := s.tokenService.ParseToken(token)
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		return false, pkg.ErrInvalidToken
	}

	playToken, err := entities.PlayTokenFromMap(claims)
	if err != nil {
		return false, pkg.ErrInvalidToken
	}

	sessionIdToken := playToken.SessionID
	sessionID := app_context.GetSessionIDFromContext(ctx)

	if sessionIdToken != sessionID {
		return false, pkg.ErrInvalidToken
	}

	sudokuDate, err := time.Parse(time.DateOnly, playToken.Date)
	if err != nil {
		return false, err
	}

	if err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		sudoku, err := s.sudokuFetcher.GetByDateAndSize(txCtx, sudokuDate, playToken.Size)
		if err != nil {
			return err
		}

		if !compareSolution(sudoku, input) {
			return pkg.ErrInvalidSolution
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}