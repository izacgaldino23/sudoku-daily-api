-- create indexes for solve table
CREATE INDEX idx_solves_sudoku_duration_user
ON solves (sudoku_id, duration, user_id);

CREATE INDEX idx_solves_duration
ON solves (duration);

CREATE INDEX idx_solves_user_duration
ON solves (user_id, duration);

CREATE INDEX idx_solves_sudoku_size_duration
ON solves (sudoku_id, size, duration);

-- create indexes for user_stats
CREATE INDEX idx_user_stats_streak_date_user
ON user_stats (longest_streak DESC, last_solved_date, user_id);

CREATE INDEX idx_user_stats_total_solved_user
ON user_stats (total_solved DESC, user_id);