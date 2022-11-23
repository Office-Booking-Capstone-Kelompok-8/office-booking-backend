package repository

import "context"

type TokenRepository interface {
	SaveToken(ctx context.Context, tokenPair string, uid string) error
}
