package mock

import (
	"context"
	"io"
	"office-booking-backend/internal/user/dto"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GetFullUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *UserServiceMock) GetAllUsers(ctx context.Context, filter *dto.UserFilterRequest) (*dto.BriefUsersResponse, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*dto.BriefUsersResponse), args.Get(1).(int64), args.Error(2)
}

func (m *UserServiceMock) GetRegisteredMemberStat(ctx context.Context) (*dto.RegisteredStatResponseList, error) {
	args := m.Called(ctx)
	return args.Get(0).(*dto.RegisteredStatResponseList), args.Error(1)
}

func (m *UserServiceMock) UpdateUserByID(ctx context.Context, id string, user *dto.UserUpdateRequest) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *UserServiceMock) DeleteUserByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *UserServiceMock) UploadUserAvatar(ctx context.Context, id string, file io.Reader) error {
	args := m.Called(ctx, id, file)
	return args.Error(0)
}
