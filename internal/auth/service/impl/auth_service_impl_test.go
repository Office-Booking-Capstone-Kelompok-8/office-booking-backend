package impl

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/internal/auth/dto"
	mockRepo "office-booking-backend/internal/auth/repository/mock"
	"office-booking-backend/internal/auth/service"
	mockToken "office-booking-backend/internal/auth/service/mock"
	mockRedis "office-booking-backend/pkg/database/redis"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	mockMail "office-booking-backend/pkg/utils/mail"
	mockPass "office-booking-backend/pkg/utils/password"
	mockRand "office-booking-backend/pkg/utils/random"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuiteAuthService struct {
	suite.Suite
	mockRepo    *mockRepo.AuthRepositoryMock
	mockToken   *mockToken.TokenServiceMock
	mockRedis   *mockRedis.RedisClientMock
	mockMail    *mockMail.ClientMock
	mockPass    *mockPass.PasswordHashService
	mockRand    *mockRand.GeneratorMock
	authService service.AuthService
}

func (s *TestSuiteAuthService) SetupTest() {
	s.mockRepo = new(mockRepo.AuthRepositoryMock)
	s.mockToken = new(mockToken.TokenServiceMock)
	s.mockRedis = new(mockRedis.RedisClientMock)
	s.mockMail = new(mockMail.ClientMock)
	s.mockPass = new(mockPass.PasswordHashService)
	s.mockRand = new(mockRand.GeneratorMock)
	s.authService = NewAuthServiceImpl(s.mockRepo, s.mockToken, s.mockRedis, s.mockMail, s.mockPass, s.mockRand)
}

func (s *TestSuiteAuthService) TearDownTest() {
	s.mockRepo = nil
	s.mockToken = nil
	s.mockPass = nil
	s.authService = nil
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(TestSuiteAuthService))
}

func (s *TestSuiteAuthService) TestRegisterUser_Success() {
	s.mockPass.On("GenerateFromPassword", []byte("password"), 10).Return([]byte("hashedPassword"), nil)
	s.mockRepo.On("RegisterUser", mock.Anything, mock.Anything).Return(nil)

	err := s.authService.RegisterUser(context.Background(), &dto.SignupRequest{
		Email:    "email",
		Password: "password",
		Name:     "name",
		Phone:    "phone",
	})

	s.NoError(err)
}

func (s *TestSuiteAuthService) TestRegisterUser_FailHashing() {
	s.mockPass.On("GenerateFromPassword", []byte("password"), 10).Return([]byte(nil), errors.New("error"))

	err := s.authService.RegisterUser(context.Background(), &dto.SignupRequest{
		Email:    "email",
		Password: "password",
		Name:     "name",
		Phone:    "phone",
	})

	s.Error(err)
}

func (s *TestSuiteAuthService) TestRegisterUser_FailRegistering() {
	s.mockPass.On("GenerateFromPassword", []byte("password"), 10).Return([]byte("hashedPassword"), nil)
	s.mockRepo.On("RegisterUser", mock.Anything, mock.Anything).Return(errors.New("error"))

	err := s.authService.RegisterUser(context.Background(), &dto.SignupRequest{
		Email:    "email",
		Password: "password",
		Name:     "name",
		Phone:    "phone",
	})

	s.Error(err)
}

func (s *TestSuiteAuthService) TestLoginUser() {
	user := &dto.LoginRequest{
		Email:    "email@mail.com",
		Password: "test",
	}
	userEntity := &entity.User{
		ID:         "someId",
		Email:      "someEmail@mail.com",
		Password:   "somePassword",
		Role:       1,
		IsVerified: false,
		Detail:     entity.UserDetail{},
	}
	tokenPair := &dto.TokenPair{
		AccessToken:  "someAccessToken",
		RefreshToken: "someRefreshToken",
	}

	for _, tc := range []struct {
		Name        string
		RepoReturn  *entity.User
		RepoError   error
		HashReturn  error
		TokenReturn *dto.TokenPair
		TokenError  error
		Expected    *dto.TokenPair
		ExpectedErr error
	}{
		{
			Name:        "Success",
			RepoReturn:  userEntity,
			RepoError:   nil,
			HashReturn:  nil,
			TokenReturn: tokenPair,
			TokenError:  nil,
			Expected:    tokenPair,
			ExpectedErr: nil,
		},
		{
			Name:        "Fail: Unknown Repo Error",
			RepoReturn:  nil,
			RepoError:   errors.New("some error"),
			Expected:    nil,
			ExpectedErr: errors.New("some error"),
		},
		{
			Name:        "Fail: User Not Found",
			RepoReturn:  nil,
			RepoError:   err2.ErrUserNotFound,
			Expected:    nil,
			ExpectedErr: err2.ErrInvalidCredentials,
		},
		{
			Name:        "Fail: Invalid Password",
			RepoReturn:  userEntity,
			RepoError:   nil,
			HashReturn:  errors.New("hash not match"),
			Expected:    nil,
			ExpectedErr: err2.ErrInvalidCredentials,
		},
		{
			Name:        "Fail: Unknown Token Error",
			RepoReturn:  userEntity,
			RepoError:   nil,
			HashReturn:  nil,
			TokenReturn: nil,
			TokenError:  errors.New("some error"),
			Expected:    nil,
			ExpectedErr: errors.New("some error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockRepo.On("GetFullUserByEmail", mock.Anything, user.Email).Return(tc.RepoReturn, tc.RepoError)
			s.mockPass.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(tc.HashReturn)
			s.mockToken.On("NewTokenPair", mock.Anything, tc.RepoReturn).Return(tc.TokenReturn, tc.TokenError)

			token, err := s.authService.LoginUser(context.Background(), user)
			s.Equal(tc.Expected, token)
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthService) TestRefreshToken() {
	for _, tc := range []struct {
		Name           string
		TokenReturn    *dto.RefreshToken
		TokenError     error
		RepoReturn     *entity.User
		RepoError      error
		GenerateReturn *dto.TokenPair
		GenerateError  error
		Expected       *dto.TokenPair
		ExpectedErr    error
	}{
		{
			Name: "Success",
			TokenReturn: &dto.RefreshToken{
				RegisteredClaims: jwt.RegisteredClaims{},
				UID:              "someId",
				Category:         "someCategory",
			},
			RepoReturn: &entity.User{
				ID:         "someId",
				IsVerified: false,
			},
			GenerateReturn: &dto.TokenPair{
				AccessToken:  "someAccessToken",
				RefreshToken: "someRefreshToken",
			},
			Expected: &dto.TokenPair{
				AccessToken:  "someAccessToken",
				RefreshToken: "someRefreshToken",
			},
		},
		{
			Name:        "Fail: Error Parsing Token",
			TokenReturn: nil,
			TokenError:  errors.New("some error"),
			Expected:    nil,
			ExpectedErr: errors.New("some error"),
		},
		{
			Name: "Fail: User Not Found",
			TokenReturn: &dto.RefreshToken{
				RegisteredClaims: jwt.RegisteredClaims{},
			},
			RepoReturn:  nil,
			RepoError:   err2.ErrUserNotFound,
			Expected:    nil,
			ExpectedErr: err2.ErrUserNotFound,
		},
		{
			Name: "Fail: Redis Error",
			TokenReturn: &dto.RefreshToken{
				RegisteredClaims: jwt.RegisteredClaims{},
			},
			RepoReturn: &entity.User{
				ID: "someId",
			},
			RepoError:      nil,
			GenerateReturn: nil,
			GenerateError:  errors.New("some error"),
			Expected:       nil,
			ExpectedErr:    errors.New("some error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockToken.On("ParseRefreshToken", mock.Anything).Return(tc.TokenReturn, tc.TokenError)
			s.mockRepo.On("GetFullUserByID", mock.Anything, mock.Anything).Return(tc.RepoReturn, tc.RepoError)
			s.mockToken.On("NewTokenPair", mock.Anything, tc.RepoReturn).Return(tc.GenerateReturn, tc.GenerateError)

			token, err := s.authService.RefreshToken(context.Background(), &dto.RefreshTokenRequest{
				RefreshToken: "someRefreshToken",
			})
			s.Equal(tc.Expected, token)
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}
