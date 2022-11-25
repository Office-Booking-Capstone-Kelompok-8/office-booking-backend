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

func (a *AuthRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := a.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (a *AuthRepositoryImpl) ChangePassword(ctx context.Context, id string, password string) error {
	result := a.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Update("password", password)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return err2.ErrUserNotFound
	}

	return nil
}
