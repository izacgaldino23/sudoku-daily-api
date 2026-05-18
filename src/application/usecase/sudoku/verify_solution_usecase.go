package sudoku

import (
	"context"
	"sudoku-daily-api/src/infrastructure/logging"
	"time"

	"sudoku-daily-api/pkg"
	user_stats_usecase "sudoku-daily-api/src/application/usecase/user_stats"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"

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
		sudoku, txErr := s.sudokuFetcher.GetByDateAndSize(txCtx, sudokuDate, playToken.Size)
		if err != nil {
			return txErr
		}

		if !compareSolution(sudoku, input) {
			return pkg.ErrInvalidSolution
		}

		solve, txErr := s.sudokuFetcher.GetSolveByIDAndUser(ctx, sudoku.ID, input.UserID)
		if txErr != nil {
			logging.Log(ctx).Info().Err(err).Msgf("error fetching solve by user %s and sudoku %s", input.UserID, sudoku.ID)
			return txErr
		}

		if solve != nil && !solve.ID.IsEmpty() && solve.Duration > 0 {
			return pkg.ErrAlreadyPlayed
		}

		txErr = s.sudokuRepo.MarkAsSolved(txCtx, solve, finished)
		if txErr != nil {
			return txErr
		}

		txErr = s.solveAddStrikeUseCase.Execute(txCtx, solve.UserID, sudokuDate)
		if txErr != nil {
			return txErr
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}
