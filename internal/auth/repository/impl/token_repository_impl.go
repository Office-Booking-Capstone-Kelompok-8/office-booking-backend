package impl

import (
	"context"
	"office-booking-backend/internal/auth/repository"

	"github.com/go-redis/redis/v9"
)

type TokenRepositoryImpl struct {
	redis *redis.Client
}

func NewTokenRepositoryImpl(redis *redis.Client) repository.TokenRepository {
	return &TokenRepositoryImpl{
		redis: redis,
	}
}

func (t *TokenRepositoryImpl) SaveToken(ctx context.Context, tokenPair string, uid string) error {
	return t.redis.Set(ctx, uid, tokenPair, 0).Err()
}
