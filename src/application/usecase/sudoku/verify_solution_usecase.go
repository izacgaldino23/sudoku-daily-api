package sudoku

import (
	"context"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/app_context"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/repository"
	"sudoku-daily-api/src/domain/vo"
	"time"
)

type (
	SudokuVerifySolutionUseCase interface {
		Execute(ctx context.Context, sudoku *entities.Solve, sessionToken string) (bool, error)
	}

	sudokuVerifySolutionUseCase struct {
		userRep       repository.UserRepository
		sudokuRepo    repository.SudokuRepository
		tokenService  domain.TokenService
		sudokuFetcher domain.SudokuDailyFetcher
	}
)

func NewSudokuVerifySolutionUseCase(
	userRep repository.UserRepository,
	sudokuRepo repository.SudokuRepository,
	tokenService domain.TokenService,
	sudokuFetcher domain.SudokuDailyFetcher,
) SudokuVerifySolutionUseCase {
	return &sudokuVerifySolutionUseCase{
		userRep:       userRep,
		sudokuRepo:    sudokuRepo,
		tokenService:  tokenService,
		sudokuFetcher: sudokuFetcher,
	}
}

func (s *sudokuVerifySolutionUseCase) Execute(ctx context.Context, solve *entities.Solve, token string) (bool, error) {
	// parse token
	claims, err := s.tokenService.ParseToken(token)
	if err != nil {
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

	// validate solution
	sudoku, err := s.sudokuFetcher.GetDaily(ctx, playToken.Size)
	if err != nil {
		return false, err
	}

	if !compareSolution(sudoku, solve) {
		return false, pkg.ErrInvalidSolution
	}

	// if logged, save on db
	if solve.UserID != "" {
		solve.ID = vo.NewUUID()
		solve.SudokuID = sudoku.ID
		solve.StartedAt = playToken.StartedAt
		solve.CompletedAt = time.Now()

		err = s.sudokuRepo.AddSolve(ctx, solve)
		if err != nil {
			return false, err
		}
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
