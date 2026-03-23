package stats

import (
	"time"

	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"

	"github.com/uptrace/bun"
)

type (
	Stats struct {
		bun.BaseModel `bun:"table:user_stats"`

		ID             string    `bun:"id,pk"`
		UserID         string    `bun:"user_id,notnull"`
		UserName       string    `bun:"username,scanonly"`
		CurrentStreak  int       `bun:",notnull"`
		LongestStreak  int       `bun:",notnull"`
		LastSolvedDate time.Time `bun:",notnull"`
		TotalSolved    int       `bun:",notnull"`
	}
)

func (s *Stats) ToDomain() *entities.UserStats {
	return &entities.UserStats{
		ID:             vo.UUID(s.ID),
		UserID:         vo.UUID(s.UserID),
		Username:       s.UserName,
		CurrentStreak:  s.CurrentStreak,
		LongestStreak:  s.LongestStreak,
		LastSolvedDate: s.LastSolvedDate,
		TotalSolved:    s.TotalSolved,
	}
}

func (s *Stats) FromDomain(stats *entities.UserStats) {
	s.ID = stats.ID.String()
	s.UserID = stats.UserID.String()
	s.CurrentStreak = stats.CurrentStreak
	s.LongestStreak = stats.LongestStreak
	s.LastSolvedDate = stats.LastSolvedDate
	s.TotalSolved = stats.TotalSolved
}
