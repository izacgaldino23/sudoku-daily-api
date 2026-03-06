package main

import (
	"fmt"
	"log"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	sudokuUsecase "sudoku-daily-api/src/application/usecase/sudoku"
	userUsecase "sudoku-daily-api/src/application/usecase/user"
	"sudoku-daily-api/src/domain/strategies"
	"sudoku-daily-api/src/infrastructure/http"
	"sudoku-daily-api/src/infrastructure/http/auth"
	"sudoku-daily-api/src/infrastructure/http/sudoku"
	"sudoku-daily-api/src/infrastructure/persistence"
	"sudoku-daily-api/src/services"

	"github.com/gofiber/fiber/v3"
)

func init() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	c := config.GetConfig()
	err = database.ConnectDB(c)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()

	configApi(app.Group("/api"))

	port := config.GetConfig().ApiPort
	fmt.Println("🚀 Server running on port", port)

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}

func configApi(app fiber.Router) {
	// others
	databaseConnection := database.GetDB()

	// strategies
	fillStrategy := strategies.NewFillStrategy()
	hideStrategy := strategies.NewHideStrategy()

	// repositories
	sudokuRepository := persistence.NewSudokuRepository(databaseConnection.BunConnection)
	userRepository := persistence.NewUserRepository(databaseConnection.BunConnection)
	refreshTokenRepository := persistence.NewRefreshTokenRepository(databaseConnection.BunConnection)

	authConfig := config.GetConfig().Auth

	// services
	generatorService := services.NewGenerator(fillStrategy, hideStrategy)
	passHasher := services.NewPasswordHasher(authConfig.Iterations, authConfig.Memory, authConfig.Parallelism, authConfig.KeyLen, authConfig.SaltLen)
	tokenService := services.NewTokenService(authConfig.SecretKey, authConfig.TokenDuration)

	// use cases
	getDailySudoku := sudokuUsecase.NewSudokuGetDailyUseCase(sudokuRepository)
	generateAll := sudokuUsecase.NewSudokuGenerateAllUseCase(sudokuRepository, generatorService)

	userRegister := userUsecase.NewUserRegisterUseCase(userRepository, passHasher)
	userLogin := userUsecase.NewUserLoginUseCase(userRepository, refreshTokenRepository, passHasher, tokenService)

	// handlers
	sudokuHandler := sudoku.NewSudokuHandler(getDailySudoku, generateAll)
	authHandler := auth.NewAuthHandler(userRegister, userLogin)

	// routes
	http.RegisterRoutes(app, sudokuHandler, authHandler)
}
