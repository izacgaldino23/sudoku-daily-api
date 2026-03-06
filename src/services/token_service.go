package services

import (
	"crypto/rand"
	"encoding/base64"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	TokenService struct{
		secret []byte
		tokenDuration int
	}
)

func NewTokenService(secret string, tokenDuration int) domain.TokenService {
	return TokenService{
		secret: []byte(secret),
		tokenDuration: tokenDuration,
	}
}

func (s TokenService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s TokenService) GenerateRefreshToken(userID string) (*entities.RefreshToken, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return nil, err
	}

	return &entities.RefreshToken{
		UserID:    vo.UUID(userID),
		Hash:      base64.URLEncoding.EncodeToString(refreshToken),
		ExpiresAt: time.Now().Add(time.Duration(s.tokenDuration) * time.Second),
	}, err
}

func (s TokenService) ValidateAccessToken(token string) (string, error) {
	panic("not implemented")
}
