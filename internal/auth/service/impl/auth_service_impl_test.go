package impl

import (
	"context"
	"errors"
	"office-booking-backend/internal/auth/dto"
	mockRepo "office-booking-backend/internal/auth/repository/mock"
	"office-booking-backend/internal/auth/service"
	mock2 "office-booking-backend/internal/auth/service/mock"

	mockPass "office-booking-backend/pkg/utils/password/mock"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuiteAuthService struct {
	suite.Suite
	mockRepo    *mockRepo.AuthRepositoryMock
	mockToken   service.TokenService
	mockPass    *mockPass.PasswordHashService
	authService service.AuthService
}

func (s *TestSuiteAuthService) SetupTest() {
	s.mockRepo = new(mockRepo.AuthRepositoryMock)
	s.mockToken = new(mock2.TokenServiceMock)
	s.mockPass = new(mockPass.PasswordHashService)
	s.authService = NewAuthServiceImpl(s.mockRepo, s.mockToken, s.mockPass)
}

func (s *TestSuiteAuthService) TearDownTest() {
	s.mockRepo = nil
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
