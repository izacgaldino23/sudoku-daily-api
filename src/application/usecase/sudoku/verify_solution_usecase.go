package sudoku

import (
	"context"
	"errors"
	"time"

	"sudoku-daily-api/pkg"
	user_stats_usecase "sudoku-daily-api/src/application/usecase/user_stats"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"

	"github.com/rs/zerolog/log"
)

type (
	SudokuVerifySolutionUseCase interface {
		Execute(ctx context.Context, sudoku *entities.Solve, playToken string, finished time.Time) (bool, error)
	}

	sudokuVerifySolutionUseCase struct {
		sudokuRepo            repository.SudokuRepository
		tokenService          domain.TokenService
		sudokuFetcher         domain.SudokuDailyFetcher
		solveAddStrikeUseCase user_stats_usecase.SolveAddStrikeUseCase
		txManager             repository.TransactionManager
	}
)

func NewSudokuVerifySolutionUseCase(
	sudokuRepo repository.SudokuRepository,
	tokenService domain.TokenService,
	sudokuFetcher domain.SudokuDailyFetcher,
	solveAddStrikeUseCase user_stats_usecase.SolveAddStrikeUseCase,
	txManager repository.TransactionManager,
) SudokuVerifySolutionUseCase {
	return &sudokuVerifySolutionUseCase{
		sudokuRepo:            sudokuRepo,
		tokenService:          tokenService,
		sudokuFetcher:         sudokuFetcher,
		solveAddStrikeUseCase: solveAddStrikeUseCase,
		txManager:             txManager,
	}
}

func (s *sudokuVerifySolutionUseCase) Execute(ctx context.Context, input *entities.Solve, token string, finished time.Time) (bool, error) {
	claims, err := s.tokenService.ParseToken(token)
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		return false, pkg.ErrInvalidToken
	}

	playToken, err := entities.PlayTokenFromMap(claims)
	if err != nil {
		return false, pkg.ErrInvalidToken
	}

	userIDFromToken := playToken.UserID
	userIDFromContext := app_context.GetUserIDFromContext(ctx)

	if userIDFromToken != userIDFromContext {
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

		solve, err := s.sudokuFetcher.GetSolveByIDAndUser(ctx, sudoku.ID, input.UserID)
		if err != nil {
			if !errors.Is(err, pkg.ErrSolutionNotFound) {
				return err
			}
		}

		if solve != nil && !solve.ID.IsEmpty() {
			return pkg.ErrAlreadyPlayed
		}
		solve = &entities.Solve{
			ID:        vo.NewUUID(),
			SudokuID:  sudoku.ID,
			Size:      sudoku.GetSize(),
			UserID:    input.UserID,
			StartedAt: playToken.StartedAt,
			Duration:  int(finished.Sub(playToken.StartedAt).Seconds()),
		}

		err = s.sudokuRepo.AddSolve(txCtx, solve)
		if err != nil {
			return err
		}

		err = s.solveAddStrikeUseCase.Execute(txCtx, solve.UserID, sudokuDate)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}