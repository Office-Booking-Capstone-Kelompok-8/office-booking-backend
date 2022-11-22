package impl

import (
	"context"
	"office-booking-backend/internal/auth/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"strings"

	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryImpl(db *gorm.DB) repository.AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (a *AuthRepositoryImpl) RegisterUser(ctx context.Context, user *entity.User) error {
	err := a.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return err2.ErrDuplicateEmail
		}

		return err
	}

	return nil
}
