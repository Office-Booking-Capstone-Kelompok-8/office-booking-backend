package impl

import (
	"context"
	"encoding/json"
	"log"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/repository"
	"office-booking-backend/internal/auth/service"
	"office-booking-backend/pkg/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenServiceImpl struct {
	AccessTokenSecret  string
	AccessTokenExp     time.Duration
	RefreshTokenSecret string
	RefreshTokenExp    time.Duration
	TokenRepository    repository.TokenRepository
}

func NewTokenServiceImpl(accessTokenSecret string, refreshTokenSecret string, accessTokenExp time.Duration, refreshTokenExp time.Duration, TokenRepository repository.TokenRepository) service.TokenService {
	return &TokenServiceImpl{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		AccessTokenExp:     accessTokenExp,
		RefreshTokenExp:    refreshTokenExp,
		TokenRepository:    TokenRepository,
	}
}

func (t *TokenServiceImpl) GenerateAccessToken(user *entity.User, exp time.Time) (string, string, error) {
	tokenID := uuid.New().String()
	accessClaims := dto.NewAccessToken(user, tokenID, exp)
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(t.AccessTokenSecret))
	if err != nil {
		log.Println("Error while signing access token:", err)
		return "", "", err
	}
	return accessToken, tokenID, nil
}

func (t *TokenServiceImpl) GenerateRefreshToken(user *entity.User, exp time.Time) (string, string, error) {
	tokenID := uuid.New().String()
	refreshClaims := dto.NewRefreshToken(user, tokenID, exp)
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(t.RefreshTokenSecret))
	if err != nil {
		log.Println("Error while signing refresh token:", err)
		return "", "", err
	}
	return refreshToken, tokenID, nil
}

func (t *TokenServiceImpl) NewTokenPair(ctx context.Context, user *entity.User) (*dto.TokenPair, error) {
	accessToken, accessID, err := t.GenerateAccessToken(user, time.Now().Add(t.AccessTokenExp))
	if err != nil {
		return nil, err
	}
	refreshToken, refreshID, err := t.GenerateRefreshToken(user, time.Now().Add(t.RefreshTokenExp))
	if err != nil {
		return nil, err
	}

	serializedTokenPair, err := json.Marshal(entity.CachedToken{
		AccessID:  accessID,
		RefreshID: refreshID,
	})
	if err != nil {
		log.Println("Error while marshalling token pair:", err)
		return nil, err
	}

	err = t.TokenRepository.SaveToken(ctx, string(serializedTokenPair), user.ID, t.RefreshTokenExp)
	if err != nil {
		log.Println("Error while saving token pair:", err)
		return nil, err
	}

	return &dto.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		nil
}
