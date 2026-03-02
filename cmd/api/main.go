package main

import (
	"fmt"
	"log"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/application/usecase"
	"sudoku-daily-api/src/infrastructure/helpers"
	"sudoku-daily-api/src/infrastructure/http"
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

	configApi(app)

	fmt.Println("🚀 Server running on port 3000")

	_ = app.Listen(config.GetConfig().ApiPort)
}

func configApi(app *fiber.App) {
	// others
	databaseConnection := database.GetDB()
	transactionManager := persistence.NewTransactionManager(databaseConnection.BunConnection)

	// helpers
	fillBacktracking := helpers.NewFillBacktracking()
	hideBacktracking := helpers.NewHideBacktracking()

	// repositories
	sudokuRepository := persistence.NewSudokuRepository(databaseConnection.BunConnection, transactionManager)

	// services
	generatorService := services.NewGenerator(fillBacktracking, hideBacktracking)

	// use cases
	getDailySudoku := usecase.NewSudokuGetDailyUseCase(sudokuRepository)
	generateAll := usecase.NewSudokuGenerateAllUseCase(sudokuRepository, generatorService)

	// handlers
	sudokuHandler := sudoku.NewSudokuHandler(getDailySudoku, generateAll)

	// routes
	http.RegisterRoutes(app, sudokuHandler)
}
