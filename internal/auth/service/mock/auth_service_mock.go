package mock

import (
	"context"
	"office-booking-backend/internal/auth/dto"

	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) RegisterUser(ctx context.Context, user *dto.SignupRequest) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *AuthServiceMock) LoginUser(ctx context.Context, user *dto.LoginRequest) (*dto.TokenPair, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*dto.TokenPair), args.Error(1)
}

func (m *AuthServiceMock) LogoutUser(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}
