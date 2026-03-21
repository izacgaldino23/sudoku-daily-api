-- create user_stats table, referencing the user, the current_streak and the longest_streak
CREATE TABLE IF NOT EXISTS user_stats (
	id UUID PRIMARY KEY,
	user_id UUID UNIQUE NOT NULL,
	current_streak INT NOT NULL DEFAULT 0,
	longest_streak INT NOT NULL DEFAULT 0,
	last_solved_date DATE NOT NULL DEFAULT CURRENT_DATE,
	total_solved INT NOT NULL DEFAULT 0,
	FOREIGN KEY (user_id) REFERENCES users(id)
)