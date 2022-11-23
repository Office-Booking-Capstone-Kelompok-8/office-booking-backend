package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type TokenRepositoryMock struct {
	mock.Mock
}

func (m *TokenRepositoryMock) SaveToken(ctx context.Context, tokenPair string, uid string, exp time.Duration) error {
	args := m.Called(ctx, tokenPair, uid, exp)
	return args.Error(0)
}

func (m *TokenRepositoryMock) DeleteToken(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *TokenRepositoryMock) GetToken(ctx context.Context, uid string) (string, error) {
	args := m.Called(ctx, uid)
	return args.String(0), args.Error(1)
}
