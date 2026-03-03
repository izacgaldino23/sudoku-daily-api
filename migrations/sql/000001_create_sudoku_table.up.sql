-- Create sudoku table
CREATE TABLE sudoku (
	id TEXT PRIMARY KEY,
	size INT NOT NULL CHECK (size in (4, 6, 9)),
	board JSONB NOT NULL,
	date TIMESTAMP NOT NULL
);