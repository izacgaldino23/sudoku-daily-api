package entities

import (
	"encoding/json"
	"sudoku-daily-api/src/domain/vo"
	"time"
)

type (
	SessionToken struct {
		Date      string    `json:"date"`
		Size      int       `json:"size"`
		StartedAt time.Time `json:"started_at"`
		SessionID vo.UUID   `json:"session_id"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	Solve struct {
		ID          vo.UUID
		SudokuID    vo.UUID
		UserID      vo.UUID
		Solution    [][]int
		StartedAt   time.Time
		CompletedAt time.Time
		CreatedAt   time.Time
	}
)

func (s *SessionToken) ToMap() map[string]any {
	return map[string]any{
		"date":       s.Date,
		"size":       s.Size,
		"started_at": s.StartedAt,
		"session_id": s.SessionID,
		"expires_at": s.ExpiresAt,
	}
}

func SessionTokenFromMap(m map[string]any) (*SessionToken, error) {
	token := &SessionToken{}

	err := json.Unmarshal([]byte(m["token"].(string)), token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
