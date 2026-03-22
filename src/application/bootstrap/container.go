package bootstrap

import (
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/repository"

	leaderboard_usecase "sudoku-daily-api/src/application/usecase/leaderboard"
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
	user_stats_usecase "sudoku-daily-api/src/application/usecase/user_stats"

	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/src/infrastructure/http/leaderboard"
	httpSudoku "sudoku-daily-api/src/infrastructure/http/sudoku"

	"github.com/gofiber/fiber/v3"
)

type Container struct {
	// infra
	DB         database.DatabaseConnection
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
	GetDailySudoku       sudokuUsecase.ISudokuGetDailyUseCase
	GenerateDailySudokus sudokuUsecase.SudokuGenerateDailyUseCase
	VerifySolution       sudokuUsecase.SudokuVerifySolutionUseCase

	UserRegister     userUsecase.UserRegisterUseCase
	UserLogin        userUsecase.UserLoginUseCase
	UserRefreshToken userUsecase.UserRefreshTokenUseCase
	UserLogout       userUsecase.UserLogoutUseCase
	UserResume       userUsecase.UserResumeUseCase

	UserStatsSolveAddStrike user_stats_usecase.SolveAddStrikeUseCase

	GetLeaderboardUseCase leaderboard_usecase.GetLeaderboard

	Middlewares Middlewares

	// handlers
	SudokuHandler      httpSudoku.SudokuHandler
	AuthHandler        auth.AuthHandler
	LeaderboardHandler leaderboard.LeaderboardHandler
}

type Middlewares struct {
	RequireJWT      fiber.Handler
	OptionalJWT     fiber.Handler
	AuthMinimum     fiber.Handler
	Session         fiber.Handler
	LogMiddleware   fiber.Handler
	RequestID       fiber.Handler
	ResponseHeaders fiber.Handler
}
