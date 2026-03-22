package entities

const (
	DailyLeaderboardType       LeaderboardType = "daily"
	AllTimeLeaderboardType     LeaderboardType = "all-time"
	TotalSolvesLeaderboardType LeaderboardType = "total"
	StreakLeaderboardType      LeaderboardType = "streak"
)

type (
	LeaderboardType string

	LeaderboardSearchParams struct {
		Type  string
		Size  BoardSize
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

func (l LeaderboardType) String() string {
	return string(l)
}
