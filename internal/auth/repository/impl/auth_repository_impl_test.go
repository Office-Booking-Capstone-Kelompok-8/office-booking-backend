package impl

import (
	"context"
	"errors"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestSuiteAuthRepository struct {
	suite.Suite
	mock sqlmock.Sqlmock
	DB   *gorm.DB
	repo *AuthRepositoryImpl
}

func (s *TestSuiteAuthRepository) SetupTest() {
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

	s.repo = &AuthRepositoryImpl{db: s.DB}
}

func (s *TestSuiteAuthRepository) TearDownTest() {
	s.mock = nil
	s.repo = nil
}

func (s *TestSuiteAuthRepository) TestImplementation() {
	s.Run("Test implementation", func() {
		s.NotPanics(func() {
			_ = NewAuthRepositoryImpl(s.DB)
		})
	})
}

func (s *TestSuiteAuthRepository) TestRegisterUser() {
	query := regexp.QuoteMeta("INSERT INTO `users` (`id`,`email`,`password`,`role`,`is_verified`,`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?,?,?,?,?,?)")
	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffected: 1,
		},
		{
			Name:         "Error: duplicate email",
			Err:          errors.New("Error 1062: Duplicate entry '' for key "),
			ExpectedErr:  err2.ErrDuplicateEmail,
			RowsAffected: 0,
		},
		{
			Name:         "Error: unknown",
			Err:          errors.New("unknown error"),
			ExpectedErr:  errors.New("unknown error"),
			RowsAffected: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			err := s.repo.RegisterUser(context.Background(), &entity.User{})

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthRepository) TestGetUserByEmail() {
	query := regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")
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

			_, err := s.repo.GetUserByEmail(context.Background(), "123")

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthRepository) TestChangePassword() {
	query := regexp.QuoteMeta("UPDATE `users` SET `password`=?,`updated_at`=? WHERE id = ? AND `users`.`deleted_at` IS NULL")
	for _, tc := range []struct {
		Name         string
		Err          error
		ExpectedErr  error
		RowsAffected int64
	}{
		{
			Name:         "Success",
			Err:          nil,
			ExpectedErr:  nil,
			RowsAffected: 1,
		},
		{
			Name:         "Error no record found",
			Err:          nil,
			ExpectedErr:  err2.ErrUserNotFound,
			RowsAffected: 0,
		},
		{
			Name:         "Error: unknown",
			Err:          errors.New("unknown error"),
			ExpectedErr:  errors.New("unknown error"),
			RowsAffected: 0,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			if tc.Err != nil {
				s.mock.ExpectExec(query).WillReturnError(tc.Err)
			} else {
				s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, tc.RowsAffected))
			}

			err := s.repo.ChangePassword(context.Background(), "123", "123")

			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func TestAuthRepository(t *testing.T) {
	suite.Run(t, new(TestSuiteAuthRepository))
}
