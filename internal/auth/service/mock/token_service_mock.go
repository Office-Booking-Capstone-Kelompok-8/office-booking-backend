package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/pkg/entity"
)

type TokenServiceMock struct {
	mock.Mock
}

func (m *TokenServiceMock) NewTokenPair(ctx context.Context, user *entity.User) (*dto.TokenPair, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*dto.TokenPair), args.Error(1)
}

func (m *TokenServiceMock) DeleteTokenPair(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}
