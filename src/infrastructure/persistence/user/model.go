package user

import (
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID            string `bun:"id,pk"`
	Username      string `bun:",unique,notnull"`
	Email         string `bun:",unique,notnull"`
	PasswordHash  string `bun:",notnull"`
	Provider      string
	ProviderID    *string `bun:",notnull"`
	EmailVerified bool    `bun:",notnull"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *User) FromDomain(user *entities.User) {
	u.ID = string(user.ID)
	u.Email = user.Email.String()
	u.Username = user.Username
	u.Provider = string(user.Provider)
	u.ProviderID = user.ProviderID
	u.EmailVerified = user.EmailVerified
	u.CreatedAt = user.CreatedAt
	u.PasswordHash = user.PasswordHash
}

func (u *User) ToDomain() *entities.User {
	return &entities.User{
		ID:            vo.UUID(u.ID),
		Email:         entities.Email(u.Email),
		Username:      u.Username,
		PasswordHash:  u.PasswordHash,
		Provider:      entities.AuthProvider(u.Provider),
		ProviderID:    u.ProviderID,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}
