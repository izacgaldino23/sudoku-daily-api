package entities

import (
	"net/mail"
	"time"

	"sudoku-daily-api/src/domain/vo"
)

const (
	EmailAuthProvider AuthProvider = "email"
)

type (
	AuthProvider string

	Email string

	User struct {
		ID            vo.UUID
		Email         Email
		Username      string
		PasswordHash  string
		Provider      AuthProvider
		ProviderID    *string
		EmailVerified bool
		CreatedAt     time.Time

		Tokens *Tokens
	}

	Tokens struct {
		AccessToken  string
		RefreshToken string
	}

	RefreshToken struct {
		ID        vo.UUID
		UserID    vo.UUID
		Hash      string
		Revoked   bool
		ExpiresAt time.Time
	}

	Resume struct {
		TotalGames map[BoardSize]int
		TodayGames []GameResult
		BestTimes  []GameResult
	}

	GameResult struct {
		Size     int
		Finished bool
		Duration int
	}
)

func (u *User) IsEmailAuth() bool {
	return u.Provider == EmailAuthProvider
}

func (e *Email) IsValid() bool {
	email, _ := mail.ParseAddress(string(*e))
	if email == nil {
		return false
	}

	*e = Email(email.Address)

	return email != nil
}

func (e *Email) String() string {
	return string(*e)
}
