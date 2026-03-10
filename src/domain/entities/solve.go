package entities

import (
	"encoding/json"
	"sudoku-daily-api/src/domain/vo"
	"time"
)

type (
	PlayToken struct {
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
		token.Size = int(size)
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
