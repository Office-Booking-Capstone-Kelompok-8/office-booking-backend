package service

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/pkg/entity"
)

type TokenService interface {
	NewTokenPair(ctx context.Context, user *entity.User) (*dto.TokenPair, error)
	DeleteTokenPair(ctx context.Context, uid string) error
	CheckToken(ctx context.Context, token *jwt.MapClaims) (bool, error)
	ParseRefreshToken(token string) (*dto.RefreshToken, error)
}
