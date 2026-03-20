package entities

const (
	DailyLeaderboardType   LeaderboardType = "daily"
	AllTimeLeaderboardType LeaderboardType = "all-time"
	StreakLeaderboardType  LeaderboardType = "streak"
)

type (
	LeaderboardType string

	LeaderboardSearchParams struct {
		Type  string
		Size  string
		Limit int
		Page  int
	}

	Leaderboard struct {
		Entries []Entry
		HasNext bool
	}

	Entry struct {
		Rank     int
		Username string
		Value    string
	}
)
