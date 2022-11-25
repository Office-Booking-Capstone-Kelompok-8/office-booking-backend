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
