package mock

import (
	"context"
	"office-booking-backend/pkg/entity"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) RegisterUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
