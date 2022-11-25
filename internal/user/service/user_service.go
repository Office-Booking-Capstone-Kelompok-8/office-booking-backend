package service

import (
	"context"
	"office-booking-backend/internal/user/dto"
)

type UserService interface {
	GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
}
