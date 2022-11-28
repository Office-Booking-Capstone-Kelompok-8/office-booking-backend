package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
	"office-booking-backend/pkg/entity"
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

func (u *UserRepositoryMock) GetAllUsers(ctx context.Context, q string, limit int, offset int) (*entity.Users, int64, error) {
	args := u.Called(ctx, q, limit, offset)
	return args.Get(0).(*entity.Users), args.Get(1).(int64), args.Error(2)
}

func (u *UserRepositoryMock) UpdateUserByID(ctx context.Context, user *entity.User) error {
	args := u.Called(ctx, user)
	return args.Error(0)
}

func (u *UserRepositoryMock) UpdateUserDetailByID(ctx context.Context, userDetail *entity.UserDetail) error {
	args := u.Called(ctx, userDetail)
	return args.Error(0)
}

func (u *UserRepositoryMock) DeleteUserByID(ctx context.Context, id string) error {
	args := u.Called(ctx, id)
	return args.Error(0)
}
