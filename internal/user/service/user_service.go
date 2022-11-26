package service

import (
	"context"
	"office-booking-backend/internal/user/dto"
)

type UserService interface {
	GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context, q string, limit int, offset int) (*dto.BriefUsersResponse, int64, error)
	UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error
}
