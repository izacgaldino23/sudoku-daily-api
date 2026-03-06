package user

import (
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID            string  `bun:"id,pk"`
	Username      string  `bun:",unique,notnull"`
	Email         string  `bun:",unique,notnull"`
	PasswordHash  []byte  `bun:",notnull"`
	Provider      string  `bun:",notnull"`
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

	if user.PasswordHash != nil {
		u.PasswordHash = []byte(*user.PasswordHash)
	}
}

func (u *User) ToDomain() *entities.User {
	var passwordHash *string
	if u.PasswordHash != nil {
		hashStr := string(u.PasswordHash)
		passwordHash = &hashStr
	}
	return &entities.User{
		ID:            vo.UUID(u.ID),
		Email:         entities.Email(u.Email),
		Username:      u.Username,
		PasswordHash:  passwordHash,
		Provider:      entities.AuthProvider(u.Provider),
		ProviderID:    u.ProviderID,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}
