package sudoku_test

import (
	"os"
	"testing"

	"sudoku-daily-api/tests/integration/helpers"
)

func TestMain(m *testing.M) {
	helpers.SetupTestEnvironment()
	defer helpers.TeardownTestDB()
	os.Exit(m.Run())
}
