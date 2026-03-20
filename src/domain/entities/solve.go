package entities

import (
	"encoding/json"
	"time"

	"sudoku-daily-api/src/domain/vo"
)

type (
	PlayToken struct {
		Date      string    `json:"date"`
		SudokuID  vo.UUID   `json:"sudoku_id"`
		SessionID vo.UUID   `json:"session_id"`
		Size      BoardSize `json:"size"`
		StartedAt time.Time `json:"started_at"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	Solve struct {
		ID        vo.UUID
		SudokuID  vo.UUID
		Size      int
		UserID    vo.UUID
		Solution  [][]int
		StartedAt time.Time
		Duration  int
		CreatedAt time.Time
	}
)

func ConvertSolvesToGameResults(solves []Solve) []GameResult {
	results := make([]GameResult, 0, len(solves))
	for _, solve := range solves {
		results = append(results, *solve.ToGameResult())
	}
	return results
}

func (s *Solve) ToGameResult() *GameResult {
	return &GameResult{
		Size:     s.Size,
		Finished: true,
		Duration: s.Duration,
	}
}

func (s *PlayToken) ToMap() map[string]any {
	return map[string]any{
		"date":       s.Date,
		"size":       s.Size,
		"started_at": s.StartedAt,
		"session_id": s.SessionID,
		"expires_at": s.ExpiresAt,
	}
}

func PlayTokenFromMap(m map[string]any) (*PlayToken, error) {
	token := &PlayToken{}

	if tokenStr, ok := m["token"].(string); ok {
		err := json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return nil, err
		}
		return token, nil
	}

	if date, ok := m["date"].(string); ok {
		token.Date = date
	}
	if size, ok := m["size"].(float64); ok {
		token.Size = BoardSize(size)
	}
	if sessionID, ok := m["session_id"].(string); ok {
		token.SessionID = vo.UUID(sessionID)
	}
	if expiresAt, ok := m["expires_at"].(string); ok {
		t, err := time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return nil, err
		}
		token.ExpiresAt = t
	}
	if startedAt, ok := m["started_at"].(string); ok {
		t, err := time.Parse(time.RFC3339, startedAt)
		if err != nil {
			return nil, err
		}
		token.StartedAt = t
	}

	return token, nil
}
