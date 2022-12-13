package repository

import (
	"context"
	"office-booking-backend/internal/user/dto"
	"office-booking-backend/pkg/entity"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) GetFullUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := u.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) GetFullUserByID(ctx context.Context, id string) (*entity.User, error) {
	args := u.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) GetAllUsers(ctx context.Context, filter *dto.UserFilterRequest) (*entity.Users, int64, error) {
	args := u.Called(ctx, filter)
	return args.Get(0).(*entity.Users), int64(args.Int(1)), args.Error(2)
}

func (u *UserRepositoryMock) GetRegisteredMemberStat(ctx context.Context) (*entity.MonthlyRegisteredStatList, error) {
	args := u.Called(ctx)
	return args.Get(0).(*entity.MonthlyRegisteredStatList), args.Error(1)
}

func (u *UserRepositoryMock) UpdateUserByID(ctx context.Context, user *entity.User) error {
	args := u.Called(ctx, user)
	return args.Error(0)
}

func (u *UserRepositoryMock) UpdateUserDetailByID(ctx context.Context, userDetail *entity.UserDetail) error {
	args := u.Called(ctx, userDetail)
	return args.Error(0)
}

func (u *UserRepositoryMock) DeleteUserByID(ctx context.Context, id string) (string, error) {
	args := u.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (u *UserRepositoryMock) GetUserProfilePictureID(ctx context.Context, id string) (*entity.ProfilePicture, error) {
	args := u.Called(ctx, id)
	return args.Get(0).(*entity.ProfilePicture), args.Error(1)
}

func (u *UserRepositoryMock) DeleteUserProfilePicture(ctx context.Context, pictureID string) error {
	args := u.Called(ctx, pictureID)
	return args.Error(0)
}
