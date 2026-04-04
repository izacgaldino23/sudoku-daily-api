package auth_test

import (
	"os"
	"testing"

	"sudoku-daily-api/tests/integration/testhelpers"
)

func TestMain(m *testing.M) {
	testhelpers.SetupTestEnvironment()
	defer testhelpers.TeardownTestDB()
	os.Exit(m.Run())
}
