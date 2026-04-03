package leaderboard

import (
	"time"

	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain/entities"
)

type (
	LeaderboardRequest struct {
		Type  string `query:"type" validate:"oneof=daily all-time streak total"`
		Size  string `query:"size"`
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

	ResetStrikesRequest struct {
		Date time.Time `json:"date" validate:"required"`
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

func (r *LeaderboardRequest) Validate() error {
	var errs pkg.ValidationErrors

	if err := pkg.ValidateStruct(r); err != nil {
		if validationErrs, ok := err.(pkg.ValidationErrors); ok {
			errs = append(errs, validationErrs...)
		}
	}

	leaderboardType := entities.LeaderboardType(r.Type)

	if leaderboardType.RequiresSize() {
		if r.Size == "" {
			errs = append(errs, pkg.ValidationError{
				Field:   "Size",
				Message: "is required for daily and all-time leaderboards",
			})
		} else if !entities.IsValidBoardSizeName(r.Size) {
			errs = append(errs, pkg.ValidationError{
				Field:   "Size",
				Message: "must be one of four, six, or nine",
			})
		}
	}

	if leaderboardType.IsSizeAllowed() && r.Size != "" {
		errs = append(errs, pkg.ValidationError{
			Field:   "Size",
			Message: "is not allowed for streak and total leaderboards",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
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
