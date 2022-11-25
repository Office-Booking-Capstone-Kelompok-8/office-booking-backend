package repository

import (
	"context"
	"office-booking-backend/internal/user/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) GetFullUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := u.db.WithContext(ctx).Model(&entity.User{}).Joins("Detail").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (u *UserRepositoryImpl) GetFullUserByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	err := u.db.WithContext(ctx).Model(&entity.User{}).Joins("Detail").Where("id = ?", id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}
