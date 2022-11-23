package repository

import (
	"context"
	"time"
)

type TokenRepository interface {
	SaveToken(ctx context.Context, tokenPair string, uid string, exp time.Duration) error
	DeleteToken(ctx context.Context, uid string) error
	GetToken(ctx context.Context, uid string) (string, error)
}
