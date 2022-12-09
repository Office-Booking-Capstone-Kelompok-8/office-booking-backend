package redis

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type RedisClientMock struct {
	mock.Mock
}

func (m *RedisClientMock) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	args := m.Called(ctx, key, value, exp)
	return args.Error(0)
}

func (m *RedisClientMock) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *RedisClientMock) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}
