-- Add solution column to store the solution of the generated Sudoku puzzles
ALTER TABLE sudokus ADD COLUMN IF NOT EXISTS solution bytea NOT NULL DEFAULT ''::bytea;