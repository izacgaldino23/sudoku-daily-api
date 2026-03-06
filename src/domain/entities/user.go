package entities

import (
	"net/mail"
	"sudoku-daily-api/src/domain/vo"
	"time"
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
		PasswordHash  *string
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
)

func (u *User) IsEmailAuth() bool {
	return u.Provider == EmailAuthProvider
}

func (e *Email) IsValid() bool {
	email, _ := mail.ParseAddress(string(*e))
	*e = Email(email.Address)

	return email != nil
}

func (e *Email) String() string {
	return string(*e)
}
