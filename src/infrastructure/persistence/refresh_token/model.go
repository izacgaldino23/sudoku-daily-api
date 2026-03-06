package refresh_token

import (
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/uptrace/bun"
)

type RefreshToken struct {
	bun.BaseModel

	ID        string    `bun:"id,pk"`
	UserID    string    `bun:"user_id,notnull"`
	TokenHash string    `bun:",unique,notnull"`
	ExpiresAt time.Time `bun:"type:timestamp,notnull"`
	Revoked   bool      `bun:"notnull"`
	CreatedAt time.Time `bun:"type:timestamp,notnull"`
}

func NewModel(token *entities.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:        vo.NewUUID().String(),
		UserID:    token.UserID.String(),
		TokenHash: token.Hash,
		ExpiresAt: token.ExpiresAt,
		Revoked:   false,
	}
}
