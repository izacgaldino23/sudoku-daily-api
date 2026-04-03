package sudoku

import (
	"context"
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
		Execute(ctx context.Context, sudoku *entities.Solve, sessionToken string, finished time.Time) (bool, error)
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

func (s *sudokuVerifySolutionUseCase) Execute(ctx context.Context, solve *entities.Solve, token string, finished time.Time) (bool, error) {
	// parse token
	claims, err := s.tokenService.ParseToken(token)
	if err != nil {
		log.Ctx(ctx).Err(err).Send()
		return false, pkg.ErrInvalidToken
	}

	playToken, err := entities.PlayTokenFromMap(claims)
	if err != nil {
		return false, pkg.ErrInvalidToken
	}

	// validate token
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
		// validate solution
		sudoku, err := s.sudokuFetcher.GetByDateAndSize(txCtx, sudokuDate, playToken.Size)
		if err != nil {
			return err
		}

		if !compareSolution(sudoku, solve) {
			return pkg.ErrInvalidSolution
		}

		// if logged, save on db
		if solve.UserID != "" {
			solve.ID = vo.NewUUID()
			solve.SudokuID = sudoku.ID
			solve.StartedAt = playToken.StartedAt
			solve.Duration = int(finished.Sub(playToken.StartedAt).Seconds())

			err = s.sudokuRepo.AddSolve(txCtx, solve)
			if err != nil {
				return err
			}

			err = s.solveAddStrikeUseCase.Execute(txCtx, solve.UserID, sudokuDate)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	return true, nil
}

func compareSolution(sudoku *entities.Sudoku, solution *entities.Solve) bool {
	board := sudoku.Solution.GetBoard()
	for i := 0; i < int(sudoku.Size); i++ {
		for j := 0; j < int(sudoku.Size); j++ {
			if board[i][j] != solution.Solution[i][j] {
				return false
			}
		}
	}

	return true
}
