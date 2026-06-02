package bootstrap

import (
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/repository"

	leaderboardUsecase "sudoku-daily-api/src/application/usecase/leaderboard"
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
	userStatsUsecase "sudoku-daily-api/src/application/usecase/user_stats"

	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/src/infrastructure/http/leaderboard"
	httpSudoku "sudoku-daily-api/src/infrastructure/http/sudoku"

	"github.com/gofiber/fiber/v3"
)

type Container struct {
	// infra
	DB         database.Connection
	Config     *config.Config
	LocalCache domain.Cache

	// repositories
	SudokuRepository       repository.SudokuRepository
	UserRepository         repository.UserRepository
	RefreshTokenRepository repository.RefreshTokenRepository
	UserStatsRepository    repository.UserStatsRepository
	TxManager              repository.TransactionManager

	// services
	GeneratorService domain.SudokuGenerator
	PasswordHasher   domain.PasswordHasher
	TokenService     domain.TokenService
	SudokuFetcher    domain.SudokuDailyFetcher
	ResumeFetcher    domain.ResumeFetcher

	// use cases
	GetDailySudoku           sudokuUsecase.ISudokuGetDailyUseCase
	GetDailySudokuForGuest   sudokuUsecase.ISudokuGetDailyForGuestUseCase
	GenerateDailySudoku      sudokuUsecase.GenerateDailyUseCase
	VerifySolution           sudokuUsecase.VerifySolutionUseCase
	VerifySolutionGuest      sudokuUsecase.VerifySolutionGuestUseCase
	GetUserSolvesUseCase     sudokuUsecase.GetUserSolvesUseCase
	RemoveUnfinishedAttempts sudokuUsecase.RemoveUnfinishedAttemptsUseCase

	UserRegister     userUsecase.RegisterUseCase
	UserLogin        userUsecase.LoginUseCase
	UserRefreshToken userUsecase.RefreshTokenUseCase
	UserLogout       userUsecase.LogoutUseCase
	UserResume       userUsecase.ResumeUseCase

	UserStatsSolveAddStrike userStatsUsecase.SolveAddStrikeUseCase

	GetLeaderboardUseCase leaderboardUsecase.GetLeaderboard
	ResetStrikesUseCase   leaderboardUsecase.ResetStrikesUseCase

	Middlewares Middlewares

	// handlers
	SudokuHandler      httpSudoku.Handler
	AuthHandler        auth.Handler
	LeaderboardHandler leaderboard.Handler
}

type Middlewares struct {
	RequireJWT        fiber.Handler
	OptionalJWT       fiber.Handler
	AuthMinimum       fiber.Handler
	Session           fiber.Handler
	LogMiddleware     fiber.Handler
	RequestID         fiber.Handler
	ResponseHeaders   fiber.Handler
	AuthOIDC          fiber.Handler
	GlobalRateLimiter fiber.Handler
	UserRateLimiter   fiber.Handler
	Timeout           fiber.Handler
	Timezone          fiber.Handler
}
