package repository

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"office-booking-backend/internal/user/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"strings"

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
	err := u.db.WithContext(ctx).
		Joins("Detail").
		Where("email = ?", email).
		First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (u *UserRepositoryImpl) GetFullUserByID(ctx context.Context, id string) (*entity.User, error) {
	db, err := u.db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := sq.Select("u.id, u.email, u.role, u.is_verified, u.created_at, u.updated_at, u.deleted_at, d.user_id, d.name, d.phone, d.picture_id, d.created_at, d.updated_at, d.deleted_at, pp.id, pp.url").
		From("users u").
		Join("user_details d ON u.id = d.user_id AND d.deleted_at IS NULL").
		LeftJoin("profile_pictures pp ON d.picture_id = pp.id").
		Where(sq.Eq{"u.deleted_at": nil, "u.id": id}).
		OrderBy("u.id").
		RunWith(db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	if !rows.Next() {
		return nil, err2.ErrUserNotFound
	}

	user := &entity.User{}
	NullAbleProfilePicture := &entity.NullAbleProfilePicture{}
	err = rows.Scan(&user.ID, &user.Email, &user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		&user.Detail.UserID, &user.Detail.Name, &user.Detail.Phone, &NullAbleProfilePicture.ID, &user.Detail.CreatedAt, &user.Detail.UpdatedAt, &user.Detail.DeletedAt,
		&NullAbleProfilePicture.ID, &NullAbleProfilePicture.Url)
	if err != nil {
		return nil, err
	}

	user.Detail.Picture = NullAbleProfilePicture.ConvertToProfilePicture()

	return user, nil
}

func (u *UserRepositoryImpl) GetAllUsers(ctx context.Context, q string, role int, limit int, offset int) (*entity.Users, int64, error) {
	users := &entity.Users{}
	var count int64

	query := u.db.WithContext(ctx).
		Joins("Detail").
		Preload("Detail.Picture")
	if q != "" {
		query = query.Where("`Detail`.`name` LIKE ?", "%"+q+"%")
	}

	err := query.
		Where("role = ?", role).
		Limit(limit).
		Offset(offset).
		Find(users).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (u *UserRepositoryImpl) UpdateUserByID(ctx context.Context, user *entity.User) error {
	res := u.db.WithContext(ctx).Where("id = ?", user.ID).Updates(user)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "for key 'users.email'") {
			return err2.ErrDuplicateEmail
		}

		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrUserNotFound
	}

	return nil
}

func (u *UserRepositoryImpl) UpdateUserDetailByID(ctx context.Context, userDetail *entity.UserDetail) error {
	res := u.db.WithContext(ctx).Where("user_id = ?", userDetail.UserID).Updates(userDetail)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrUserNotFound
	}

	return nil
}

func (u *UserRepositoryImpl) DeleteUserByID(ctx context.Context, id string) (string, error) {
	// res := u.db.WithContext(ctx).Delete(&entity.User{ID: id})
	// if res.Error != nil {
	// 	return res.Error
	// }
	//
	// if res.RowsAffected == 0 {
	// 	return err2.ErrUserNotFound
	// }
	//
	// return nil

	var pictureId sql.NullString
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.UserDetail{}).
			Select("picture_id").
			Where("user_id = ?", id).
			First(&pictureId).Error
		isHaveProfilePicture := true
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}

			isHaveProfilePicture = false
		}

		err = tx.Unscoped().Delete(&entity.UserDetail{UserID: id}).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Delete(&entity.User{ID: id}).Error
		if err != nil {
			return err
		}

		if isHaveProfilePicture {
			err = tx.Unscoped().Where("id = ?", pictureId.String).Delete(&entity.ProfilePicture{}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	return pictureId.String, err
}

func (u *UserRepositoryImpl) GetUserProfilePictureID(ctx context.Context, id string) (*entity.ProfilePicture, error) {
	userDetail := &entity.UserDetail{}
	err := u.db.WithContext(ctx).
		Model(&entity.UserDetail{}).
		Select("picture_id").
		Joins("Picture").
		Where("user_id = ?", id).
		First(userDetail).Error
	if err != nil {
		return nil, err
	}

	return &userDetail.Picture, nil
}

func (u *UserRepositoryImpl) DeleteUserProfilePicture(ctx context.Context, pictureID string) error {
	res := u.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", pictureID).
		Delete(&entity.ProfilePicture{})
	if res.Error != nil {
		return res.Error
	}

	return nil
}
