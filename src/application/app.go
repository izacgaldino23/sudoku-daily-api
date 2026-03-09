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
	"sudoku-daily-api/src/infrastructure/persistence/cache"
	persistenceRefreshToken "sudoku-daily-api/src/infrastructure/persistence/database/refresh_token"
	persistenceSudoku "sudoku-daily-api/src/infrastructure/persistence/database/sudoku"
	persistenceTx "sudoku-daily-api/src/infrastructure/persistence/database/tx"
	persistenceUser "sudoku-daily-api/src/infrastructure/persistence/database/user"
	"sudoku-daily-api/src/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/rs/zerolog/log"
)

const (
	maxCacheSize = 5
)

func InitApp(app fiber.Router) error {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// storages
	databaseConnection := database.GetDB()
	localCache := cache.NewLocalCache(maxCacheSize)

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
	getDailySudoku := sudokuUsecase.NewSudokuGetDailyUseCase(sudokuRepository, localCache)
	generateAll := sudokuUsecase.NewSudokuGenerateAllUseCase(txManager, sudokuRepository, generatorService)

	userRegister := userUsecase.NewUserRegisterUseCase(userRepository, passHasher)
	userLogin := userUsecase.NewUserLoginUseCase(txManager, userRepository, refreshTokenRepository, passHasher, tokenService)
	userRefreshToken := userUsecase.NewUserRefreshTokenUseCase(refreshTokenRepository, tokenService)
	userLogoutUseCase := userUsecase.NewUserLogoutUseCase(refreshTokenRepository)

	// middlewares
	tokenMiddleware := middlewares.JWTMiddleware(tokenService)
	authMiddleware := middlewares.AuthMiddleware(tokenService)
	logMiddleware := middlewares.LogMiddleware(log.Logger)

	// handlers
	sudokuHandler := httpSudoku.NewSudokuHandler(getDailySudoku, generateAll)
	authHandler := auth.NewAuthHandler(userRegister, userLogin, userRefreshToken, userLogoutUseCase)

	app.Use(logMiddleware)

	// routes
	http.RegisterRoutes(app, sudokuHandler, authHandler, tokenMiddleware, authMiddleware)

	return nil
}
