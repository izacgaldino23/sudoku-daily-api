package main

import (
	"log"
	"sudoku-daily-api/migrations"
	"sudoku-daily-api/pkg/config"
	"sudoku-daily-api/pkg/database"
)

func init() {
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err = database.ConnectDB(c); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := migrations.RunMigrations(); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations executed successfully")
}
