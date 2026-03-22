package leaderboard

import "sudoku-daily-api/src/domain/entities"

type (
	LeaderboardRequest struct {
		Type  string `query:"type" validate:"oneof=daily all-time streak"`
		Size  string `query:"size" validate:"oneof=four six nine"`
		Limit int    `query:"limit" validate:"min=1,max=100"`
		Page  int    `query:"page" validate:"min=1"`
	}

	LeaderboardResponse struct {
		Entries []Entry `json:"solves"`
		HasNext bool    `json:"has_next"`
	}

	Entry struct {
		Rank     int    `json:"rank"`
		Username string `json:"username"`
		Value    string `json:"value"`
	}
)

func (r *LeaderboardRequest) ToDomain() *entities.LeaderboardSearchParams {
	return &entities.LeaderboardSearchParams{
		Type:  r.Type,
		Size:  entities.BoardSizeFromName(r.Size),
		Limit: r.Limit,
		Page:  r.Page,
	}
}

func responseFromDomain(leaderboard *entities.Leaderboard) LeaderboardResponse {
	var entries []Entry

	for _, entry := range leaderboard.Entries {
		entries = append(entries, Entry{
			Rank:     entry.Rank,
			Username: entry.Username,
			Value:    entry.Value,
		})
	}

	return LeaderboardResponse{
		Entries: entries,
		HasNext: leaderboard.HasNext,
	}
}
