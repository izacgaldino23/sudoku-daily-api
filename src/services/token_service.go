package services

import (
	"crypto/rand"
	"encoding/base64"
	"sudoku-daily-api/pkg"
	"sudoku-daily-api/src/domain"
	"sudoku-daily-api/src/domain/entities"
	"sudoku-daily-api/src/domain/vo"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	TokenService struct {
		secret               []byte
		accessTokenDuration  int
		refreshTokenDuration int
	}
)

func NewTokenService(
	secret string,
	accessTokenDuration int,
	refreshTokenDuration int) domain.TokenService {
	return &TokenService{
		secret:               []byte(secret),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *TokenService) GenerateJWTToken(fields map[string]any) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Duration(s.accessTokenDuration) * time.Second).Unix(),
		"iat": time.Now().Unix(),
	}

	for key, value := range fields {
		claims[key] = value
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenService) GenerateRefreshToken(userID vo.UUID) (*entities.RefreshToken, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return nil, err
	}

	return &entities.RefreshToken{
		UserID:    userID,
		Hash:      base64.URLEncoding.EncodeToString(refreshToken),
		ExpiresAt: time.Now().Add(time.Duration(s.refreshTokenDuration) * time.Second),
	}, err
}

func (s *TokenService) ValidateAccessToken(token string) (vo.UUID, error) {
	claims, err := s.ParseToken(token)
	if err != nil {
		return "", err
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return "", pkg.ErrInvalidToken
	}

	if claims["iat"].(float64) > float64(time.Now().Unix()) {
		return "", pkg.ErrInvalidToken
	}

	if claims["user_id"] == nil {
		return "", pkg.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", pkg.ErrInvalidToken
	}

	return vo.UUID(userID), nil
}

func (s *TokenService) ParseToken(token string) (map[string]any, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
