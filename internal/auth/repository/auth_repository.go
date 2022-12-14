package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ChangePassword(ctx context.Context, id string, password string) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	VerifyEmail(ctx context.Context, id string) error
}
