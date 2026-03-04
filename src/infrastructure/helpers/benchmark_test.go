package helpers

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"testing"
	"time"
)

func BenchmarkFillBacktracking4(b *testing.B) {
	benchmarkFillBacktracking(b, entities.BoardSize4)
}

func BenchmarkFillBacktracking6(b *testing.B) {
	benchmarkFillBacktracking(b, entities.BoardSize6)
}

func BenchmarkFillBacktracking9(b *testing.B) {
	benchmarkFillBacktracking(b, entities.BoardSize9)
}

func benchmarkFillBacktracking(b *testing.B, size entities.BoardSize) {
	f := NewFillBacktracking()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sudoku := entities.NewSudoku(size)
		f.Fill(sudoku, r)
	}
}

func BenchmarkHideBacktracking4(b *testing.B) {
	benchmarkHideBacktracking(b, entities.BoardSize4)
}

func BenchmarkHideBacktracking6(b *testing.B) {
	benchmarkHideBacktracking(b, entities.BoardSize6)
}

func BenchmarkHideBacktracking9(b *testing.B) {
	benchmarkHideBacktracking(b, entities.BoardSize9)
}

func benchmarkHideBacktracking(b *testing.B, size entities.BoardSize) {
	h := NewHideBacktracking()
	f := NewFillBacktracking()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sudoku := entities.NewSudoku(size)
		sudoku.Difficulty = entities.DifficultyMedium
		f.Fill(sudoku, r)
		h.Hide(sudoku, r)
	}
}

func BenchmarkSolver4(b *testing.B) {
	benchmarkSolver(b, entities.BoardSize4)
}

func BenchmarkSolver6(b *testing.B) {
	benchmarkSolver(b, entities.BoardSize6)
}

func BenchmarkSolver9(b *testing.B) {
	benchmarkSolver(b, entities.BoardSize9)
}

func benchmarkSolver(b *testing.B, size entities.BoardSize) {
	solver := NewSolver()
	f := NewFillBacktracking()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sudoku := entities.NewSudoku(size)
	sudoku.Difficulty = entities.DifficultyMedium
	f.Fill(sudoku, r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.Execute(sudoku)
	}
}

func BenchmarkSolverEmpty4(b *testing.B) {
	benchmarkSolverEmpty(b, entities.BoardSize4)
}

func BenchmarkSolverEmpty6(b *testing.B) {
	benchmarkSolverEmpty(b, entities.BoardSize6)
}

func BenchmarkSolverEmpty9(b *testing.B) {
	benchmarkSolverEmpty(b, entities.BoardSize9)
}

func benchmarkSolverEmpty(b *testing.B, size entities.BoardSize) {
	solver := NewSolver()
	f := NewFillBacktracking()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sudoku := entities.NewSudoku(size)
	sudoku.Difficulty = entities.DifficultyMedium
	f.Fill(sudoku, r)

	emptyBoard := entities.NewSudoku(size)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for r := 0; r < int(size); r++ {
			for c := 0; c < int(size); c++ {
				emptyBoard.Board.SetCell(r, c, sudoku.Board.GetCell(r, c))
			}
		}
		solver.Execute(emptyBoard)
	}
}
