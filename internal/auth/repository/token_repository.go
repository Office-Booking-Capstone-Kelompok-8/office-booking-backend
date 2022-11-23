package repository

import (
	"context"
	"time"
)

type TokenRepository interface {
	SaveToken(ctx context.Context, tokenPair string, uid string, exp time.Duration) error
}
