package entities

import (
	"net/mail"
	"time"
)

const (
	EmailAuthProvider AuthProvider = "email"
)

type (
	AuthProvider string

	UserID string

	Email string

	User struct {
		ID            UserID
		Email         Email
		Username      string
		PasswordHash  *string
		Provider      AuthProvider
		ProviderID    *string
		EmailVerified bool
		CreatedAt     time.Time
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

func (u *UserID) String() string {
	return string(*u)
}
