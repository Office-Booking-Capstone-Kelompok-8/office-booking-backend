package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type UserRepository interface {
	GetFullUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetFullUserByID(ctx context.Context, id string) (*entity.User, error)
	GetAllUsers(ctx context.Context, q string, limit int, offset int) (*entity.Users, int64, error)
	UpdateUserByID(ctx context.Context, user *entity.User) error
	UpdateUserDetailByID(ctx context.Context, userDetail *entity.UserDetail) error
}
