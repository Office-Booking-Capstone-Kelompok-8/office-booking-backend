package service

import (
	"context"
	"io"
	"office-booking-backend/internal/user/dto"
)

type UserService interface {
	GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context, q string, role int, limit int, offset int) (*dto.BriefUsersResponse, int64, error)
	GetRegisteredMemberStat(ctx context.Context) (*dto.RegisteredStatResponseList, error)
	UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error
	DeleteUserByID(ctx context.Context, id string) error
	UploadUserAvatar(ctx context.Context, id string, file io.Reader) error
}
