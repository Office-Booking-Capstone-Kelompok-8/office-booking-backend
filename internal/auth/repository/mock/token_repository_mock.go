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
