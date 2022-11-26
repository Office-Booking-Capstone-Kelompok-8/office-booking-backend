package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"office-booking-backend/internal/user/dto"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *UserServiceMock) GetAllUsers(ctx context.Context, q string, limit int, offset int) (*dto.BriefUsersResponse, int64, error) {
	args := m.Called(ctx, q, limit, offset)
	return args.Get(0).(*dto.BriefUsersResponse), args.Get(1).(int64), args.Error(2)
}

func (m *UserServiceMock) UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}
