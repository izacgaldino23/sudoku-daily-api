-- Add solution column to store the solution of the generated Sudoku puzzles
ALTER TABLE sudoku ADD COLUMN IF NOT EXISTS solution bytea NOT NULL DEFAULT ''::bytea;