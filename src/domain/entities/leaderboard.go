package entities

const (
	DailyLeaderboardType       LeaderboardType = "daily"
	AllTimeLeaderboardType     LeaderboardType = "all-time"
	TotalSolvesLeaderboardType LeaderboardType = "total"
	StreakLeaderboardType      LeaderboardType = "streak"
)

var (
	ValidLeaderboardTypes = []LeaderboardType{
		DailyLeaderboardType,
		AllTimeLeaderboardType,
		TotalSolvesLeaderboardType,
		StreakLeaderboardType,
	}
)

func (l LeaderboardType) IsValid() bool {
	for _, valid := range ValidLeaderboardTypes {
		if l == valid {
			return true
		}
	}
	return false
}

func (l LeaderboardType) RequiresSize() bool {
	return l == DailyLeaderboardType || l == AllTimeLeaderboardType
}

func (l LeaderboardType) IsSizeAllowed() bool {
	return !l.RequiresSize()
}

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
