-- Add solution column to store the solution of the generated Sudoku puzzles
ALTER TABLE sudoku ADD COLUMN solution bytea NOT NULL;