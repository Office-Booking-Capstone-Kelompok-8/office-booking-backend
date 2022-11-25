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
