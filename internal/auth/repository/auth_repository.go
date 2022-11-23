package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, user *entity.User) error
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
