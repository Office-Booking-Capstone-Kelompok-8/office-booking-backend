package service

import (
	"context"
	"office-booking-backend/internal/auth/dto"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *dto.SignupRequest) error
}
