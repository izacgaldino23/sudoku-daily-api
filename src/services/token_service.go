package services

import (
	"crypto/rand"
	"encoding/base64"
	"sudoku-daily-api/src/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	TokenService struct{
		secret []byte
	}
)

func NewTokenService(secret string) domain.TokenService {
	return TokenService{
		secret: []byte(secret),
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

func (s TokenService) GenerateRefreshToken(userID string) (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(refreshToken), err
}

func (s TokenService) ValidateAccessToken(token string) (string, error) {
	panic("not implemented")
}
