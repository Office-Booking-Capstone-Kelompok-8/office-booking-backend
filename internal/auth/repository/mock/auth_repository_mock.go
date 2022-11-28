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

func (m *AuthRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *AuthRepositoryMock) ChangePassword(ctx context.Context, id string, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}
