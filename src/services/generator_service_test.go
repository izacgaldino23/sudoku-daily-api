package services

// func TestGenerator(t *testing.T) {
// 	t.Run("TestGenerator_isLineValid", func(t *testing.T) {
// 		g := &sudokuGenerator{}

// 		validCases := [][][]int{
// 			{
// 				{1, 2, 3, 4},
// 				{3, 4, 1, 2},
// 				{2, 3, 4, 1},
// 				{4, 1, 2, 3},
// 			},
// 		}

// 		testIsLineValid(t, g, validCases, true)

// 		invalidCases := [][][]int{
// 			{
// 				{1, 2, 3, 3},
// 				{2, 4, 4, 1},
// 				{1, 4, 1, 2},
// 				{4, 1, 3, 3},
// 			},
// 		}

// 		testIsLineValid(t, g, invalidCases, false)
// 	})

// 	t.Run("TestGenerator_generateTiles", func(t *testing.T) {
// 		g := &sudokuGenerator{}

// 		size := 6
// 		sudoku := g.GenerateDaily(size, 1000)

// 		for _, row := range sudoku.Board {
// 			// each row should have unique numbers from 1 to size
// 			for i := 1; i <= size; i++ {
// 				assert.Contains(t, row, i)
// 			}

// 			// each row sum should be equal to size * (size + 1) / 2
// 			var sum int
// 			for _, v := range row {
// 				sum += v
// 			}
// 			assert.Equal(t, sum, size*(size+1)/2)

// 			t.Log(row)
// 		}
// 	})
// }

// func testIsLineValid(t *testing.T, g *sudokuGenerator, cases [][][]int, valid bool) {
// 	for _, board := range cases {
// 		size := len(board)
// 		// lines
// 		for i := 0; i < size; i++ {
// 			r := g.isLineValid(board, i, 0, 1, size)
// 			assert.Equal(t, valid, r, "line expected %v", valid)
// 		}
// 		// column
// 		for i := 0; i < size; i++ {
// 			r := g.isLineValid(board, 0, i, size, 1)
// 			assert.Equal(t, valid, r, "column expected %v", valid)
// 		}

// 		// subGrids
// 		grids := g.getGrids(size)
// 		for _, grid := range grids {
// 			r := g.isLineValid(board, grid.row, grid.col, grid.rowCount, grid.colCount)
// 			assert.Equal(t, valid, r, "grid %v", grid)
// 		}
// 	}
// }
