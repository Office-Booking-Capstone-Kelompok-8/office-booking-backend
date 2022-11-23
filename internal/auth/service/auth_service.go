package service

import (
	"context"
	"office-booking-backend/internal/auth/dto"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *dto.SignupRequest) error
	LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error)
	LogoutUser(ctx context.Context, uid string) error
}
