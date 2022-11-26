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

func (a *UserRepositoryImpl) GetAllUsers(ctx context.Context, q string, limit int, offset int) (*entity.Users, int64, error) {
	users := &entity.Users{}
	var count int64

	query := a.db.WithContext(ctx).Joins("Detail").Model(&entity.User{})
	if q != "" {
		query = query.Where("`Detail`.`name` LIKE ?", "%"+q+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(users).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (u *UserRepositoryImpl) UpdateUserByID(ctx context.Context, user *entity.User) error {
	res := u.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrUserNotFound
	}

	return nil
}

func (u *UserRepositoryImpl) UpdateUserDetailByID(ctx context.Context, userDetail *entity.UserDetail) error {
	res := u.db.WithContext(ctx).Model(&entity.UserDetail{}).Where("user_id = ?", userDetail.UserID).Updates(userDetail)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrUserNotFound
	}

	return nil
}
