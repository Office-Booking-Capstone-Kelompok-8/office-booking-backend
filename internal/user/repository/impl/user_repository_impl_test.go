package repository

import (
	"context"
	"errors"
	"office-booking-backend/internal/user/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestSuiteUserRepository struct {
	suite.Suite
	mock sqlmock.Sqlmock
	DB   *gorm.DB
	repo *UserRepositoryImpl
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteUserRepository))
}

func (s *TestSuiteUserRepository) SetupTest() {
	mockConn, mock, err := sqlmock.New()
	s.Require().NoError(err)

	s.mock = mock
	s.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockConn,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	s.Require().NoError(err)

	s.repo = &UserRepositoryImpl{db: s.DB}
}

func (s *TestSuiteUserRepository) TearDownTest() {
	s.mock = nil
	s.repo = nil
}

func (s *TestSuiteUserRepository) TestNewUserRepositoryImpl() {
	s.Run("Success", func() {
		repo := NewUserRepositoryImpl(s.DB)
		s.Implements((*repository.UserRepository)(nil), repo)
	})
}

func (s *TestSuiteUserRepository) TestGetFullUserByEmail() {
	query := regexp.QuoteMeta("SELECT `users`.`id`,`users`.`email`,`users`.`password`,`users`.`role`,`users`.`is_verified`,`users`.`created_at`,`users`.`updated_at`,`users`.`deleted_at`,`Detail`.`user_id` AS `Detail__user_id`,`Detail`.`name` AS `Detail__name`,`Detail`.`phone` AS `Detail__phone`,`Detail`.`picture_id` AS `Detail__picture_id`,`Detail`.`created_at` AS `Detail__created_at`,`Detail`.`updated_at` AS `Detail__updated_at`,`Detail`.`deleted_at` AS `Detail__deleted_at` FROM `users` LEFT JOIN `user_details` `Detail` ON `users`.`id` = `Detail`.`user_id` AND `Detail`.`deleted_at` IS NULL WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")
	for _, tc := range []struct {
		Name        string
		Err         error
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Error: no record found",
			Err:         gorm.ErrRecordNotFound,
			ExpectedErr: err2.ErrUserNotFound,
		},
		{
			Name:        "Error: unknown",
			Err:         errors.New("unknown error"),
			ExpectedErr: errors.New("unknown error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "role"}).AddRow(1, "123", "123", 1))
			}

			_, err := s.repo.GetFullUserByEmail(context.Background(), "123")

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserRepository) TestGetFullUserByID() {
	header := sqlmock.NewRows([]string{"id", "username", "role", "is_verified", "created_at", "updated_at", "deleted_at", "user_id", "name", "phone", "picture_id", "created_at", "updated_at", "deleted_at", "id", "url"})
	query := regexp.QuoteMeta("SELECT u.id,u.email,u.role,u.is_verified,u.created_at,u.updated_at,u.deleted_at,d.user_id ,d.name,d.phone ,d.picture_id ,d.created_at ,d.updated_at ,d.deleted_at,pp.id, pp.url FROM users u JOIN user_details d ON u.id = d.user_id AND d.deleted_at IS NULL LEFT JOIN profile_pictures pp ON d.picture_id = pp.id WHERE u.deleted_at IS NULL AND u.id = ? ORDER BY u.id")
	for _, tc := range []struct {
		Name        string
		Rows        *sqlmock.Rows
		Err         error
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Rows:        header.AddRow(1, "123", 1, true, time.Now(), time.Now(), time.Now(), 1, "123", "123", 1, time.Now(), time.Now(), time.Now(), 1, "123"),
			Err:         nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Error: no record found",
			Rows:        header,
			ExpectedErr: err2.ErrUserNotFound,
		},
		{
			Name:        "Error: unknown",
			Err:         errors.New("unknown error"),
			ExpectedErr: errors.New("unknown error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(tc.Rows)
			}

			_, err := s.repo.GetFullUserByID(context.Background(), "123")

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserRepository) TestGetAllUsers() {
	query := regexp.QuoteMeta("SELECT `users`.`id`,`users`.`email`,`users`.`password`,`users`.`role`,`users`.`is_verified`,`users`.`created_at`,`users`.`updated_at`,`users`.`deleted_at`,`Detail`.`user_id` AS `Detail__user_id`,`Detail`.`name` AS `Detail__name`,`Detail`.`phone` AS `Detail__phone`,`Detail`.`picture_id` AS `Detail__picture_id`,`Detail`.`created_at` AS `Detail__created_at`,`Detail`.`updated_at` AS `Detail__updated_at`,`Detail`.`deleted_at` AS `Detail__deleted_at` FROM `users` LEFT JOIN `user_details` `Detail` ON `users`.`id` = `Detail`.`user_id` AND `Detail`.`deleted_at` IS NULL WHERE `Detail`.`name` LIKE ? AND role = ? AND `users`.`deleted_at` IS NULL LIMIT 1 OFFSET 1")
	count := regexp.QuoteMeta("SELECT count(*) FROM `users` LEFT JOIN `user_details` `Detail` ON `users`.`id` = `Detail`.`user_id` AND `Detail`.`deleted_at` IS NULL WHERE `Detail`.`name` LIKE ? AND role = ? AND `users`.`deleted_at` IS NULL LIMIT 1 OFFSET 1")
	for _, tc := range []struct {
		Name        string
		Err         error
		ExpectedErr error
	}{
		{
			Name:        "Success",
			Err:         nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Error: unknown",
			Err:         errors.New("unknown error"),
			ExpectedErr: errors.New("unknown error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectQuery(query).WillReturnError(tc.Err)
				s.mock.ExpectQuery(count).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "role"}).AddRow(1, "123", "123", 1))
				s.mock.ExpectQuery(count).WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
			}

			_, _, err := s.repo.GetAllUsers(context.Background(), "123", 1, 1, 1)

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserRepository) TestUpdateUserByID() {
	query := regexp.QuoteMeta("UPDATE `users` SET `updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")
	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffrcted int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffrcted: 1,
		},
		{
			Name:         "Error: no record found",
			Err:          nil,
			ExpectedErr:  err2.ErrUserNotFound,
			RowsAffrcted: 0,
		},
		{
			Name:         "Error: duplicate entry",
			Err:          errors.New("Error 1062: Duplicate entry '123' for key 'users.email'"),
			ExpectedErr:  err2.ErrDuplicateEmail,
			RowsAffrcted: 0,
		},
		{
			Name:         "Error: unknown",
			Err:          errors.New("unknown error"),
			ExpectedErr:  errors.New("unknown error"),
			RowsAffrcted: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffrcted))
			}

			err := s.repo.UpdateUserByID(context.Background(), &entity.User{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserRepository) TestUpdateUserDetailByID() {
	query := regexp.QuoteMeta("UPDATE `user_details` SET `updated_at`=? WHERE user_id = ? AND `user_details`.`deleted_at` IS NULL")
	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffrcted int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffrcted: 1,
		},
		{
			Name:         "Error: no record found",
			Err:          nil,
			ExpectedErr:  err2.ErrUserNotFound,
			RowsAffrcted: 0,
		},
		{
			Name:         "Error: unknown",
			Err:          errors.New("unknown error"),
			ExpectedErr:  errors.New("unknown error"),
			RowsAffrcted: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffrcted))
			}

			err := s.repo.UpdateUserDetailByID(context.Background(), &entity.UserDetail{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}
