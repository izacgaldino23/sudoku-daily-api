package auth

import (
	"sudoku-daily-api/src/domain/entities"
	"time"
)

type (
	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`

		UserName  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	RefreshTokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	LogoutRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)

func (r *RegisterRequest) ToDomain() *entities.User {
	return &entities.User{
		Email:        entities.Email(r.Email),
		Username:     r.Username,
		PasswordHash: &r.Password,
	}
}

func (r *LoginRequest) ToDomain() *entities.User {
	return &entities.User{
		Email:        entities.Email(r.Email),
		PasswordHash: &r.Password,
	}
}

func (r *LoginResponse) FromDomain(user *entities.User) {
	r.UserName = user.Username
	r.Email = string(user.Email)
	r.CreatedAt = user.CreatedAt

	if user.Tokens != nil {
		r.AccessToken = user.Tokens.AccessToken
		r.RefreshToken = user.Tokens.RefreshToken
	}
}

func (r *RefreshTokenResponse) FromDomain(accessToken string) {
	r.AccessToken = accessToken
}
