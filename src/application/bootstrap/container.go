package bootstrap

import (
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/repository"

	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"

	"sudoku-daily-api/src/infrastructure/http/auth"
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
	TxManager              repository.TransactionManager

	// services
	GeneratorService domain.SudokuGenerator
	PasswordHasher   domain.PasswordHasher
	TokenService     domain.TokenService
	SudokuFetcher    domain.SudokuDailyFetcher

	// use cases
	GetDailySudoku     sudokuUsecase.ISudokuGetDailyUseCase
	GenerateAllSudokus sudokuUsecase.SudokuGenerateAllUseCase
	VerifySolution     sudokuUsecase.SudokuVerifySolutionUseCase

	UserRegister     userUsecase.UserRegisterUseCase
	UserLogin        userUsecase.UserLoginUseCase
	UserRefreshToken userUsecase.UserRefreshTokenUseCase
	UserLogout       userUsecase.UserLogoutUseCase

	// middlewares
	RequireJWT          fiber.Handler
	OptionalJWT         fiber.Handler
	AuthMinimum         fiber.Handler
	Session             fiber.Handler
	LogMiddleware       fiber.Handler
	RequestIDMiddleware fiber.Handler

	// handlers
	SudokuHandler httpSudoku.SudokuHandler
	AuthHandler   auth.AuthHandler
}
