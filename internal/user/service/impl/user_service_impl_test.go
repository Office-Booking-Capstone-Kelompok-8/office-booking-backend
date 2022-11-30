package impl

import (
	"context"
	mockReservationSrv "office-booking-backend/internal/reservation/service/mock"
	"office-booking-backend/internal/user/dto"
	mockRepo "office-booking-backend/internal/user/repository/mock"
	"office-booking-backend/internal/user/service"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	mockImageKitSrv "office-booking-backend/pkg/utils/imagekit"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuiteUserService struct {
	suite.Suite
	mockRepo           *mockRepo.UserRepositoryMock
	mockReservationSrv *mockReservationSrv.ReservationServiceMock
	mockImageKitSrv    *mockImageKitSrv.ImgKitServiceMock
	userService        service.UserService
}

func (s *TestSuiteUserService) SetupTest() {
	s.mockRepo = new(mockRepo.UserRepositoryMock)
	s.mockReservationSrv = new(mockReservationSrv.ReservationServiceMock)
	s.mockImageKitSrv = new(mockImageKitSrv.ImgKitServiceMock)
	s.userService = NewUserServiceImpl(s.mockRepo, s.mockReservationSrv, s.mockImageKitSrv)
}

func (s *TestSuiteUserService) TearDownTest() {
	s.mockRepo = nil
	s.mockReservationSrv = nil
	s.mockImageKitSrv = nil
	s.userService = nil
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestSuiteUserService))
}
func (s *TestSuiteUserService) TestGetFullUserByID_Success() {
	s.mockRepo.On("GetFullUserByID", mock.Anything, mock.Anything).Return(&entity.User{}, nil)
	_, err := s.userService.GetFullUserByID(context.Background(), "")
	s.NoError(err)
}

func (s *TestSuiteUserService) TestGetFullUserByID_Fail() {
	s.mockRepo.On("GetFullUserByID", mock.Anything, mock.Anything).Return((*entity.User)(nil), err2.ErrUserNotFound)
	_, err := s.userService.GetFullUserByID(context.Background(), "")
	s.Error(err)
}

func (s *TestSuiteUserService) TestGetAllUsers_Success() {
	s.mockRepo.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&entity.Users{}, int64(0), nil)
	_, _, err := s.userService.GetAllUsers(context.Background(), "", 1, 0, 0)
	s.NoError(err)
}

func (s *TestSuiteUserService) TestGetAllUsers_Fail() {
	s.mockRepo.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*entity.Users)(nil), int64(0), err2.ErrUserNotFound)
	_, _, err := s.userService.GetAllUsers(context.Background(), "", 1, 0, 0)
	s.Error(err)
}

func (s *TestSuiteUserService) TestUpdateUserByID() {
	for _, tc := range []struct {
		Name          string
		User          *dto.UserUpdateRequest
		UserEntitty   *entity.User
		UserRepoErr   error
		DetailRepoErr error
		ExpectedErr   error
	}{
		{
			Name:          "Success",
			User:          &dto.UserUpdateRequest{},
			UserEntitty:   &entity.User{},
			UserRepoErr:   nil,
			DetailRepoErr: nil,
			ExpectedErr:   nil,
		},
		{
			Name: "Success: with changed email",
			User: &dto.UserUpdateRequest{
				Email: "123@mail.com",
			},
			UserEntitty: &entity.User{
				Email:      "123@mail.com",
				IsVerified: false,
			},
			UserRepoErr:   nil,
			DetailRepoErr: nil,
			ExpectedErr:   nil,
		},
		{
			Name:          "Fail: User update error",
			User:          &dto.UserUpdateRequest{},
			UserEntitty:   &entity.User{},
			UserRepoErr:   err2.ErrUserNotFound,
			DetailRepoErr: nil,
			ExpectedErr:   err2.ErrUserNotFound,
		},
		{
			Name:          "Fail: User detail update error",
			User:          &dto.UserUpdateRequest{},
			UserEntitty:   &entity.User{},
			UserRepoErr:   nil,
			DetailRepoErr: err2.ErrUserNotFound,
			ExpectedErr:   err2.ErrUserNotFound,
		},
		{
			Name:          "Fail: User and user detail update error",
			User:          &dto.UserUpdateRequest{},
			UserEntitty:   &entity.User{},
			UserRepoErr:   err2.ErrUserNotFound,
			DetailRepoErr: err2.ErrUserNotFound,
			ExpectedErr:   err2.ErrUserNotFound,
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockRepo.On("UpdateUserByID", mock.Anything, tc.UserEntitty).Return(tc.UserRepoErr)
			s.mockRepo.On("UpdateUserDetailByID", mock.Anything, &tc.UserEntitty.Detail).Return(tc.DetailRepoErr)
			err := s.userService.UpdateUserByID(context.Background(), "", tc.User)
			s.Equal(tc.ExpectedErr, err)
		})
		s.TearDownTest()
	}
}
