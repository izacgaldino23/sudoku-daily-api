package auth

import "sudoku-daily-api/src/domain/entities"

type (
	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=3,max=100"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}
)

func (r *RegisterRequest) ToDomain() *entities.User {
	return &entities.User{
		Email:        entities.Email(r.Email),
		Username:     r.Username,
		PasswordHash: &r.Password,
	}
}
