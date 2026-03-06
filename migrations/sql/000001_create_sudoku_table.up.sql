-- Create sudoku table
CREATE TABLE IF NOT EXISTS sudoku (
	id UUID PRIMARY KEY,
	size INT NOT NULL CHECK (size in (4, 6, 9)),
	difficulty TEXT NOT NULL CHECK (difficulty in ('easy', 'medium', 'hard')),
	board bytea NOT NULL,
	date DATE NOT NULL
);