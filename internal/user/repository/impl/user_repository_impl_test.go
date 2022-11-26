package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	err2 "office-booking-backend/pkg/errors"
	"regexp"
	"testing"
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
	query := regexp.QuoteMeta("SELECT `users`.`id`,`users`.`email`,`users`.`password`,`users`.`role`,`users`.`is_verified`,`users`.`created_at`,`users`.`updated_at`,`users`.`deleted_at`,`Detail`.`user_id` AS `Detail__user_id`,`Detail`.`name` AS `Detail__name`,`Detail`.`phone` AS `Detail__phone`,`Detail`.`picture_id` AS `Detail__picture_id`,`Detail`.`created_at` AS `Detail__created_at`,`Detail`.`updated_at` AS `Detail__updated_at`,`Detail`.`deleted_at` AS `Detail__deleted_at` FROM `users` LEFT JOIN `user_details` `Detail` ON `users`.`id` = `Detail`.`user_id` AND `Detail`.`deleted_at` IS NULL WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")
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

			_, err := s.repo.GetFullUserByID(context.Background(), "123")

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserRepository) TestGetAllUsers() {
	query := regexp.QuoteMeta("SELECT `users`.`id`,`users`.`email`,`users`.`password`,`users`.`role`,`users`.`is_verified`,`users`.`created_at`,`users`.`updated_at`,`users`.`deleted_at`,`Detail`.`user_id` AS `Detail__user_id`,`Detail`.`name` AS `Detail__name`,`Detail`.`phone` AS `Detail__phone`,`Detail`.`profile_picture_id` AS `Detail__profile_picture_id`,`Detail`.`created_at` AS `Detail__created_at`,`Detail`.`updated_at` AS `Detail__updated_at`,`Detail`.`deleted_at` AS `Detail__deleted_at` FROM `users` LEFT JOIN `user_details` `Detail` ON `users`.`id` = `Detail`.`user_id` AND `Detail`.`deleted_at` IS NULL WHERE `Detail`.`name` LIKE ? AND `users`.`deleted_at` IS NULL LIMIT 1 OFFSET 1")
	count := regexp.QuoteMeta("SELECT count(*) FROM `users` LEFT JOIN `user_details` `Detail` ON `users`.`id` = `Detail`.`user_id` AND `Detail`.`deleted_at` IS NULL WHERE `Detail`.`name` LIKE ? AND `users`.`deleted_at` IS NULL LIMIT 1 OFFSET 1")
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

			_, _, err := s.repo.GetAllUsers(context.Background(), "123", 1, 1)

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}
