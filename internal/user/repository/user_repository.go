package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type UserRepository interface {
	GetFullUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetFullUserByID(ctx context.Context, id string) (*entity.User, error)
}
