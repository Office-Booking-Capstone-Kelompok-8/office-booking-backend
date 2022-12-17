package service

import (
	"context"
	"io"
	"office-booking-backend/internal/user/dto"
)

type UserService interface {
	GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context, filter *dto.UserFilterRequest) (*dto.BriefUsersResponse, int64, error)
	GetRegisteredMemberStat(ctx context.Context) (*dto.RegisteredStatResponseList, error)
	GetRegisteredMemberCount(ctx context.Context) (*dto.TotalByTimeFrame, error)
	UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error
	DeleteUserByID(ctx context.Context, id string) error
	UploadUserAvatar(ctx context.Context, id string, file io.Reader) error
}
