package service

import (
	"context"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/pkg/entity"
)

type TokenService interface {
	NewTokenPair(ctx context.Context, user *entity.User) (*dto.TokenPair, error)
	DeleteTokenPair(ctx context.Context, uid string) error
}
