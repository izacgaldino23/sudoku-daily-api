package strategies

import (
	"math/rand"
	"sudoku-daily-api/src/domain/entities"
	"testing"
	"time"
)

func BenchmarkFillStrategy4(b *testing.B) {
	benchmarkFillStrategy(b, entities.BoardSize4)
}

func BenchmarkFillStrategy6(b *testing.B) {
	benchmarkFillStrategy(b, entities.BoardSize6)
}

func BenchmarkFillStrategy9(b *testing.B) {
	benchmarkFillStrategy(b, entities.BoardSize9)
}

func benchmarkFillStrategy(b *testing.B, size entities.BoardSize) {
	f := NewFillStrategy()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sudoku := entities.NewSudoku(size)
		f.Fill(sudoku, r)
	}
}

func BenchmarkHideStrategy4(b *testing.B) {
	benchmarkHideStrategy(b, entities.BoardSize4)
}

func BenchmarkHideStrategy6(b *testing.B) {
	benchmarkHideStrategy(b, entities.BoardSize6)
}

func BenchmarkHideStrategy9(b *testing.B) {
	benchmarkHideStrategy(b, entities.BoardSize9)
}

func benchmarkHideStrategy(b *testing.B, size entities.BoardSize) {
	h := NewHideStrategy()
	f := NewFillStrategy()
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
	solver := newSolver()
	f := NewFillStrategy()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sudoku := entities.NewSudoku(size)
	sudoku.Difficulty = entities.DifficultyMedium
	f.Fill(sudoku, r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.Execute(&sudoku.Board)
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
	solver := newSolver()
	f := NewFillStrategy()
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
		solver.Execute(&emptyBoard.Board)
	}
}
