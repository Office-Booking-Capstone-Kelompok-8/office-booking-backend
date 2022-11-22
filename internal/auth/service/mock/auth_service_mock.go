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
