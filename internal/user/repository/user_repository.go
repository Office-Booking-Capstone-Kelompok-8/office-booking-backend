package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type UserRepository interface {
	GetFullUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetFullUserByID(ctx context.Context, id string) (*entity.User, error)
	GetAllUsers(ctx context.Context, q string, role int, limit int, offset int) (*entity.Users, int64, error)
	GetUserProfilePictureID(ctx context.Context, id string) (*entity.ProfilePicture, error)
	GetRegisteredMemberStat(ctx context.Context) (*entity.MonthlyRegisteredStatList, error)
	UpdateUserByID(ctx context.Context, user *entity.User) error
	UpdateUserDetailByID(ctx context.Context, userDetail *entity.UserDetail) error
	DeleteUserByID(ctx context.Context, id string) (string, error)
	DeleteUserProfilePicture(ctx context.Context, pictureID string) error
}
