package impl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/internal/auth/dto"
	mockRepo "office-booking-backend/internal/auth/repository/mock"
	mockToken "office-booking-backend/internal/auth/service/mock"
	mockRepo2 "office-booking-backend/internal/user/repository/mock"
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
	mockRepo     *mockRepo.AuthRepositoryMock
	mockUserRepo *mockRepo2.UserRepositoryMock
	mockToken    *mockToken.TokenServiceMock
	mockRedis    *mockRedis.RedisClientMock
	mockMail     *mockMail.ClientMock
	mockPass     *mockPass.PasswordHashService
	mockRand     *mockRand.GeneratorMock
	authService  *AuthServiceImpl
}

func (s *TestSuiteAuthService) SetupTest() {
	s.mockRepo = new(mockRepo.AuthRepositoryMock)
	s.mockUserRepo = new(mockRepo2.UserRepositoryMock)
	s.mockToken = new(mockToken.TokenServiceMock)
	s.mockRedis = new(mockRedis.RedisClientMock)
	s.mockMail = new(mockMail.ClientMock)
	s.mockPass = new(mockPass.PasswordHashService)
	s.mockRand = new(mockRand.GeneratorMock)
	s.authService = &AuthServiceImpl{
		repository: s.mockRepo,
		userRepo:   s.mockUserRepo,
		token:      s.mockToken,
		redisRepo:  s.mockRedis,
		mail:       s.mockMail,
		password:   s.mockPass,
		generator:  s.mockRand,
	}
}

func (s *TestSuiteAuthService) TearDownTest() {
	s.mockRepo = nil
	s.mockUserRepo = nil
	s.mockToken = nil
	s.mockRedis = nil
	s.mockMail = nil
	s.mockPass = nil
	s.mockRand = nil
	s.authService = nil
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(TestSuiteAuthService))
}

func (s *TestSuiteAuthService) TestImplementation() {
	s.Run("Success", func() {
		s.NotPanics(func() {
			_ = NewAuthServiceImpl(s.mockRepo, s.mockUserRepo, s.mockToken, s.mockRedis, s.mockMail, s.mockPass, s.mockRand)
		})
	})
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
			s.mockUserRepo.On("GetFullUserByEmail", mock.Anything, user.Email).Return(tc.RepoReturn, tc.RepoError)
			s.mockPass.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(tc.HashReturn)
			s.mockToken.On("NewTokenPair", mock.Anything, tc.RepoReturn).Return(tc.TokenReturn, tc.TokenError)

			token, err := s.authService.LoginUser(context.Background(), user)
			s.Equal(tc.Expected, token)
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthService) TestLogoutUser() {
	s.Run("Success", func() {
		s.mockToken.On("DeleteTokenPair", mock.Anything, mock.Anything).Return(nil)
		err := s.authService.LogoutUser(context.Background(), "someToken")
		s.NoError(err)
	})
}

func (s *TestSuiteAuthService) TestCreateKey() {
	s.Run("Success", func() {
		s.NotPanics(func() {
			key := createKey("someText")
			s.NotEmpty(key)
		})
	})
}

func (s *TestSuiteAuthService) TestCreateOTP() {
	for _, tc := range []struct {
		Name          string
		RepoReturn    *entity.User
		RepoError     error
		GenerateError error
		RedisErr      error
		ExpectedErr   error
	}{
		{
			Name:          "Success",
			RepoReturn:    &entity.User{},
			RepoError:     nil,
			GenerateError: nil,
			RedisErr:      nil,
			ExpectedErr:   nil,
		},
		{
			Name:          "Fail: Unknown Repo Error",
			RepoReturn:    nil,
			RepoError:     errors.New("some error"),
			GenerateError: nil,
			RedisErr:      nil,
			ExpectedErr:   errors.New("some error"),
		},
		{
			Name:          "Fail: User Not Found",
			RepoReturn:    nil,
			RepoError:     err2.ErrUserNotFound,
			GenerateError: nil,
			RedisErr:      nil,
			ExpectedErr:   err2.ErrUserNotFound,
		},
		{
			Name:          "Fail: Unknown Generate Error",
			RepoReturn:    &entity.User{},
			RepoError:     nil,
			GenerateError: errors.New("some error"),
			RedisErr:      nil,
			ExpectedErr:   errors.New("some error"),
		},
		{
			Name:          "Fail: Unknown Redis Error",
			RepoReturn:    &entity.User{},
			RepoError:     nil,
			GenerateError: nil,
			RedisErr:      errors.New("some error"),
			ExpectedErr:   errors.New("some error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(tc.RepoReturn, tc.RepoError)
			s.mockRand.On("GenerateRandomIntString", 6).Return("123456", tc.GenerateError)
			s.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.RedisErr)

			_, err := s.authService.createOTP(context.Background(), "someEmail")
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthService) TestRequestOTP() {
	for _, tc := range []struct {
		Name        string
		MailErr     error
		ExpectedErr error
	}{
		{
			Name:        "Success",
			MailErr:     nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Fail: Unknown Mail Error",
			MailErr:     errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
			s.mockRand.On("GenerateRandomIntString", 6).Return("123456", nil)
			s.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			s.mockMail.On("SendMail", mock.Anything, mock.Anything).Return(tc.MailErr)
			err := s.authService.RequestOTP(context.Background(), "someEmail")
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthService) TestVerifyOTP() {
	for _, tc := range []struct {
		Name        string
		otp         *dto.OTPVerifyRequest
		RedisReturn entity.CachedOTP
		RedisErr    error
		ExpectedErr error
	}{
		{
			Name: "Success",
			otp: &dto.OTPVerifyRequest{
				Email: "someEmail",
				Code:  "123456",
			},
			RedisReturn: entity.CachedOTP{
				OTP: "123456",
			},
		},
		{
			Name: "Fail: OTP not found",
			otp: &dto.OTPVerifyRequest{
				Email: "someEmail",
				Code:  "123456",
			},
			RedisErr:    redis.Nil,
			ExpectedErr: err2.ErrInvalidOTP,
		},
		{
			Name: "Fail: OTP not match",
			otp: &dto.OTPVerifyRequest{
				Email: "someEmail",
				Code:  "123456",
			},
			RedisReturn: entity.CachedOTP{
				OTP: "1234567",
			},
			ExpectedErr: err2.ErrInvalidOTP,
		},
		{
			Name: "Fail: Unknown Redis Error",
			otp: &dto.OTPVerifyRequest{
				Email: "someEmail",
				Code:  "123456",
			},
			RedisErr:    errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonRedisReturn, err := json.Marshal(tc.RedisReturn)
			s.NoError(err)

			s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(string(jsonRedisReturn), tc.RedisErr)
			_, err = s.authService.VerifyOTP(context.Background(), tc.otp)
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
			s.mockUserRepo.On("GetFullUserByID", mock.Anything, mock.Anything).Return(tc.RepoReturn, tc.RepoError)
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

func (s *TestSuiteAuthService) TestResetPassword() {
	for _, tc := range []struct {
		Name        string
		RedisReturn entity.CachedOTP
		RedisErr    error
		Redis2Err   error
		HashErr     error
		RepoErr     error
		ExpectedErr error
	}{
		{
			Name: "Success",
			RedisReturn: entity.CachedOTP{
				Key: "someKey",
			},
			RedisErr:    nil,
			Redis2Err:   nil,
			HashErr:     nil,
			ExpectedErr: nil,
		},
		{
			Name:        "Fail: Redis Get OTP Error",
			RedisErr:    errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
		{
			Name:        "Fail: Redis Get OTP Not Found",
			RedisErr:    redis.Nil,
			ExpectedErr: err2.ErrInvalidOTPToken,
		},
		{
			Name: "Fail: Key Not Match",
			RedisReturn: entity.CachedOTP{
				Key: "someOtherKey",
			},
			ExpectedErr: err2.ErrInvalidOTPToken,
		},
		{
			Name: "Fail: Redis Delete OTP Error",
			RedisReturn: entity.CachedOTP{
				Key: "someKey",
			},
			RedisErr:    nil,
			Redis2Err:   errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
		{
			Name: "Fail: Hashing Error",
			RedisReturn: entity.CachedOTP{
				Key: "someKey",
			},
			RedisErr:    nil,
			Redis2Err:   nil,
			HashErr:     errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
		{
			Name: "Fail: Repo Error",
			RedisReturn: entity.CachedOTP{
				Key: "someKey",
			},
			RedisErr:    nil,
			Redis2Err:   nil,
			HashErr:     nil,
			RepoErr:     errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
		{
			Name: "Fail: User Not Found",
			RedisReturn: entity.CachedOTP{
				Key: "someKey",
			},
			RedisErr:    nil,
			Redis2Err:   nil,
			HashErr:     nil,
			RepoErr:     err2.ErrUserNotFound,
			ExpectedErr: err2.ErrUserNotFound,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonRedisReturn, err := json.Marshal(tc.RedisReturn)
			s.NoError(err)

			s.mockRedis.On("Get", mock.Anything, mock.Anything).Return(string(jsonRedisReturn), tc.RedisErr)
			s.mockRedis.On("Del", mock.Anything, mock.Anything).Return(tc.Redis2Err)
			s.mockPass.On("GenerateFromPassword", mock.Anything, mock.Anything).Return([]byte("someHash"), tc.HashErr)
			s.mockRepo.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(tc.RepoErr)

			err = s.authService.ResetPassword(context.Background(), &dto.PasswordResetRequest{
				Email:    "someEmail",
				Password: "123",
				Key:      "someKey",
			})
			s.Equal(tc.ExpectedErr, err)

		})
		s.TearDownTest()
	}
}

func (s *TestSuiteAuthService) TestChangePassword() {
	for _, tc := range []struct {
		Name           string
		GetRepoErr     error
		CompareHashErr error
		HashErr        error
		ChangeErr      error
		ExpectedErr    error
	}{
		{
			Name: "Success",
		},
		{
			Name:        "Fail: Get Repo Error",
			GetRepoErr:  errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
		{
			Name:        "Fail: User Not Found",
			GetRepoErr:  err2.ErrUserNotFound,
			ExpectedErr: err2.ErrUserNotFound,
		},
		{
			Name:           "Fail: Compare Hash Error or Password Not Match",
			CompareHashErr: errors.New("some error"),
			ExpectedErr:    err2.ErrPasswordNotMatch,
		},
		{
			Name:        "Fail: Error on Hashing New Password",
			HashErr:     errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
		{
			Name:        "Fail: Error on Change Password",
			ChangeErr:   errors.New("some error"),
			ExpectedErr: errors.New("some error"),
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockRepo.On("GetUserByID", mock.Anything, mock.Anything).Return(&entity.User{}, tc.GetRepoErr)
			s.mockPass.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(tc.CompareHashErr)
			s.mockPass.On("GenerateFromPassword", mock.Anything, mock.Anything).Return([]byte("someHash"), tc.HashErr)
			s.mockRepo.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(tc.ChangeErr)

			err := s.authService.ChangePassword(context.Background(), "someUid", &dto.ChangePasswordRequest{
				OldPassword: "123",
				NewPassword: "1234",
			})
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}
