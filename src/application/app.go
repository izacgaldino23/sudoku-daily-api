package application

import (
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
	"sudoku-daily-api/src/domain/strategies"
	"sudoku-daily-api/src/infrastructure/http"
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/src/infrastructure/http/middlewares"
	httpSudoku "sudoku-daily-api/src/infrastructure/http/sudoku"
	persistenceRefreshToken "sudoku-daily-api/src/infrastructure/persistence/refresh_token"
	persistenceSudoku "sudoku-daily-api/src/infrastructure/persistence/sudoku"
	persistenceTx "sudoku-daily-api/src/infrastructure/persistence/tx"
	persistenceUser "sudoku-daily-api/src/infrastructure/persistence/user"
	"sudoku-daily-api/src/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog/log"
)

func InitApp(app fiber.Router) error {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// others
	databaseConnection := database.GetDB()

	// strategies
	fillStrategy := strategies.NewFillStrategy()
	hideStrategy := strategies.NewHideStrategy()

	// repositories
	sudokuRepository := persistenceSudoku.NewRepository(databaseConnection.BunConnection)
	userRepository := persistenceUser.NewRepository(databaseConnection.BunConnection)
	refreshTokenRepository := persistenceRefreshToken.NewRepository(databaseConnection.BunConnection)

	txManager := persistenceTx.NewTransactionManager(databaseConnection.BunConnection)

	authConfig := config.GetConfig().Auth

	// services
	generatorService := services.NewGenerator(fillStrategy, hideStrategy)
	passHasher := services.NewPasswordHasher(authConfig.Iterations, authConfig.Memory, authConfig.Parallelism, authConfig.KeyLen, authConfig.SaltLen)
	tokenService := services.NewTokenService(authConfig.SecretKey, authConfig.AccessTokenDuration, authConfig.RefreshTokenDuration)

	// use cases
	getDailySudoku := sudokuUsecase.NewSudokuGetDailyUseCase(sudokuRepository)
	generateAll := sudokuUsecase.NewSudokuGenerateAllUseCase(txManager, sudokuRepository, generatorService)

	userRegister := userUsecase.NewUserRegisterUseCase(userRepository, passHasher)
	userLogin := userUsecase.NewUserLoginUseCase(txManager, userRepository, refreshTokenRepository, passHasher, tokenService)
	userRefreshToken := userUsecase.NewUserRefreshTokenUseCase(refreshTokenRepository, tokenService)
	userLogoutUseCase := userUsecase.NewUserLogoutUseCase(refreshTokenRepository)

	// middlewares
	authMiddleware := middlewares.NewJWTMiddleware(tokenService)
	logMiddleware := middlewares.NewLogMiddleware()

	// handlers
	sudokuHandler := httpSudoku.NewSudokuHandler(getDailySudoku, generateAll)
	authHandler := auth.NewAuthHandler(userRegister, userLogin, userRefreshToken, userLogoutUseCase)

	app.Use(logMiddleware.Execute(log.Logger))

	// routes
	http.RegisterRoutes(app, authMiddleware, sudokuHandler, authHandler)

	return nil
}
