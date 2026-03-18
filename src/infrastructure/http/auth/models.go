package auth

import (
	"time"

	"sudoku-daily-api/src/domain/entities"
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

	ResumeResponse struct {
		TotalGames map[int]int  `json:"total_games"`
		TodayGames []gameResult `json:"today_games"`
		BestTimes  []gameResult `json:"best_times"`
	}

	gameResult struct {
		Size     int  `json:"size"`
		Finished bool `json:"finished"`
		Time     int  `json:"time"`
	}
)

func (r *RegisterRequest) ToDomain() *entities.User {
	return &entities.User{
		Email:        entities.Email(r.Email),
		Username:     r.Username,
		PasswordHash: r.Password,
		Provider:     entities.EmailAuthProvider,
	}
}

func (r *LoginRequest) ToDomain() *entities.User {
	return &entities.User{
		Email:        entities.Email(r.Email),
		PasswordHash: r.Password,
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

func (r *ResumeResponse) FromDomain(resume *entities.Resume) {
	if resume == nil {
		return
	}

	r.BestTimes = make([]gameResult, len(resume.BestTimes))
	for i := range resume.BestTimes {
		r.BestTimes[i] = gameResult{
			Size:     resume.BestTimes[i].Size,
			Finished: resume.BestTimes[i].Finished,
			Time:     resume.BestTimes[i].Time,
		}
	}

	r.TodayGames = make([]gameResult, len(resume.TodayGames))
	for i := range resume.TodayGames {
		r.TodayGames[i] = gameResult{
			Size:     resume.TodayGames[i].Size,
			Finished: resume.TodayGames[i].Finished,
			Time:     resume.TodayGames[i].Time,
		}
	}

	r.TotalGames = make(map[int]int, len(resume.TotalGames))
	for size, count := range resume.TotalGames {
		r.TotalGames[int(size)] = count
	}
}