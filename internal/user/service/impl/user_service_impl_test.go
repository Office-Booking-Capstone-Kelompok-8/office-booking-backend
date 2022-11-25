package impl

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	mockRepo "office-booking-backend/internal/user/repository/mock"
	"office-booking-backend/internal/user/service"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"testing"
)

type TestSuiteUserService struct {
	suite.Suite
	mockRepo    *mockRepo.UserRepositoryMock
	userService service.UserService
}

func (s *TestSuiteUserService) SetupTest() {
	s.mockRepo = new(mockRepo.UserRepositoryMock)
	s.userService = NewUserServiceImpl(s.mockRepo)
}

func (s *TestSuiteUserService) TearDownTest() {
	s.mockRepo = nil
	s.userService = nil
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestSuiteUserService))
}
func (s *TestSuiteUserService) TestGetFullUserByID_Success() {
	s.mockRepo.On("GetFullUserByID", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
	_, err := s.userService.GetFullUserByID(nil, "")
	s.NoError(err)
}

func (s *TestSuiteUserService) TestGetFullUserByID_Fail() {
	s.mockRepo.On("GetFullUserByID", mock.Anything, mock.Anything).Return((*entity.User)(nil), err2.ErrUserNotFound)
	_, err := s.userService.GetFullUserByID(nil, "")
	s.Error(err)
}
