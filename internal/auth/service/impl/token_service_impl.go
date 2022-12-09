package impl

import (
	"context"
	"encoding/json"
	"log"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/service"
	"office-booking-backend/pkg/database/redis"
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
	RedisRepository    redis.RedisClient
}

func NewTokenServiceImpl(accessTokenSecret string, refreshTokenSecret string, accessTokenExp time.Duration, refreshTokenExp time.Duration, RedisRepository redis.RedisClient) service.TokenService {
	return &TokenServiceImpl{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		AccessTokenExp:     accessTokenExp,
		RefreshTokenExp:    refreshTokenExp,
		RedisRepository:    RedisRepository,
	}
}

func (t *TokenServiceImpl) generateAccessToken(user *entity.User, exp time.Time) (string, string, error) {
	tokenID := uuid.New().String()
	accessClaims := dto.NewAccessToken(user, tokenID, exp)
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(t.AccessTokenSecret))
	if err != nil {
		log.Println("Error while signing access token:", err)
		return "", "", err
	}
	return accessToken, tokenID, nil
}

func (t *TokenServiceImpl) generateRefreshToken(user *entity.User, exp time.Time) (string, string, error) {
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
	accessToken, accessID, err := t.generateAccessToken(user, time.Now().Add(t.AccessTokenExp))
	if err != nil {
		return nil, err
	}
	refreshToken, refreshID, err := t.generateRefreshToken(user, time.Now().Add(t.RefreshTokenExp))
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

	err = t.RedisRepository.Set(ctx, user.ID, string(serializedTokenPair), t.RefreshTokenExp)
	if err != nil {
		log.Println("Error while saving token pair:", err)
		return nil, err
	}

	return &dto.TokenPair{
			Role:         user.Role,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		nil
}

func (t *TokenServiceImpl) DeleteTokenPair(ctx context.Context, uid string) error {
	err := t.RedisRepository.Del(ctx, uid)
	if err != nil {
		log.Println("Error while deleting token pair:", err)
		return err
	}
	return nil
}

func (t *TokenServiceImpl) CheckToken(ctx context.Context, token *jwt.MapClaims) (bool, error) {
	cachedJson, err := t.RedisRepository.Get(ctx, (*token)["uid"].(string))
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}

		log.Println("Error while getting token pair:", err)
		return false, err
	}

	cachedToken := new(entity.CachedToken)
	err = json.Unmarshal([]byte(cachedJson), cachedToken)

	var cachedID string
	if (*token)["cat"].(string) == "access" {
		cachedID = cachedToken.AccessID
	} else {
		cachedID = cachedToken.RefreshID
	}

	if err != nil || cachedID != (*token)["jti"].(string) {
		return false, nil
	}

	return true, nil
}

func (t *TokenServiceImpl) ParseRefreshToken(tokenStr string) (*dto.RefreshToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &dto.RefreshToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.RefreshTokenSecret), nil
	})
	if err != nil {
		log.Println("Error while parsing refresh token:", err)
		return nil, err
	}

	return token.Claims.(*dto.RefreshToken), nil
}
