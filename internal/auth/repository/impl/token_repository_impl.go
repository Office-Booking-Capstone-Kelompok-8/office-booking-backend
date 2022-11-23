package impl

import (
	"context"
	"office-booking-backend/internal/auth/repository"
	"time"

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

func (t *TokenRepositoryImpl) SaveToken(ctx context.Context, tokenPair string, uid string, exp time.Duration) error {
	return t.redis.Set(ctx, uid, tokenPair, exp).Err()
}

func (t *TokenRepositoryImpl) DeleteToken(ctx context.Context, uid string) error {
	return t.redis.Del(ctx, uid).Err()
}

func (t *TokenRepositoryImpl) GetToken(ctx context.Context, uid string) (string, error) {
	return t.redis.Get(ctx, uid).Result()
}
