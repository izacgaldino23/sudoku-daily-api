-- create solve table referencing the user and the sudoku
CREATE TABLE IF NOT EXISTS solves (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	sudoku_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	started_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (sudoku_id) REFERENCES sudokus(id) ON DELETE CASCADE
);