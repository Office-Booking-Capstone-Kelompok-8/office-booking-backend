package dto

import (
	"office-booking-backend/pkg/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessToken struct {
	jwt.RegisteredClaims
	UID        string `json:"uid"`
	Role       int    `json:"role"`
	IsVerified bool   `json:"isVerified"`
	Category   string `json:"cat"`
}

func NewAccessToken(user *entity.User, tokenID string, exp time.Time) *AccessToken {
	return &AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		UID:        user.ID,
		Role:       user.Role,
		IsVerified: *user.IsVerified,
		Category:   "access",
	}
}

type RefreshToken struct {
	jwt.RegisteredClaims
	UID      string `json:"uid"`
	Category string `json:"cat"`
}

func NewRefreshToken(user *entity.User, tokenID string, exp time.Time) *RefreshToken {
	return &RefreshToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		UID:      user.ID,
		Category: "refresh",
	}
}

type TokenPair struct {
	Role         int    `json:"role"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
