package main

import (
	"fmt"
	"log"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
	"sudoku-daily-api/src/application"

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

	apiRouter := app.Group("/api")
	application.InitApp(apiRouter)

	port := config.GetConfig().ApiPort
	fmt.Println("🚀 Server running on port", port)

	err := app.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}
