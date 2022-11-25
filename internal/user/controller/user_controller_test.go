package controller

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"office-booking-backend/internal/user/dto"
	mockService "office-booking-backend/internal/user/service/mock"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"testing"
)

type TestSuiteUserController struct {
	suite.Suite
	mockService    *mockService.UserServiceMock
	userController *UserController
	fiberApp       *fiber.App
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(TestSuiteUserController))
}

func (s *TestSuiteUserController) SetupTest() {
	s.mockService = new(mockService.UserServiceMock)
	s.userController = NewUserController(s.mockService)
	s.fiberApp = fiber.New(fiber.Config{
		ErrorHandler: response.DefaultErrorHandler,
	})
}

func (s *TestSuiteUserController) TearDownTest() {
	s.mockService = nil
	s.userController = nil
	s.fiberApp = nil
}

func (s *TestSuiteUserController) TestGetLoggedFullUserByID() {
	token := &jwt.Token{
		Claims: jwt.MapClaims{
			"uid": "some_uid",
		},
	}

	for _, tc := range []struct {
		Name           string
		ServiceReturns *dto.UserResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
		ExpectedErr    error
	}{
		{
			Name: "success",
			ServiceReturns: &dto.UserResponse{
				ID: "some_uid",
			},
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user fetched successfully",
				Data: map[string]interface{}{
					"email":   "",
					"id":      "some_uid",
					"name":    "",
					"phone":   "",
					"picture": "",
					"role":    float64(0),
				},
			},
		},
		{
			Name:           "Fail: User not found",
			ServiceReturns: nil,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusNotFound,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name:           "Fail: Unknown error",
			ServiceReturns: nil,
			ServiceErr:     errors.New("some error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "some error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockService.On("GetFullUserByID", mock.Anything, mock.Anything).Return(tc.ServiceReturns, tc.ServiceErr)

			s.fiberApp.Post("/", func(ctx *fiber.Ctx) error {
				ctx.Locals("user", token) // set token to context
				return s.userController.GetLoggedFullUserByID(ctx)
			})
			r := httptest.NewRequest("POST", "/", nil)
			resp, err := s.fiberApp.Test(r)
			s.NoError(err)

			var body response.BaseResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			s.NoError(err)

			s.Equal(tc.ExpectedStatus, resp.StatusCode)
			s.Equal(tc.ExpectedBody, body)
		})
		s.TearDownTest()
	}
}

func (s *TestSuiteUserController) TestGetFullUserByID() {
	token := &jwt.Token{
		Claims: jwt.MapClaims{
			"uid": "some_uid",
		},
	}

	for _, tc := range []struct {
		Name           string
		ServiceReturns *dto.UserResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
		ExpectedErr    error
	}{
		{
			Name: "success",
			ServiceReturns: &dto.UserResponse{
				ID: "some_uid",
			},
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user fetched successfully",
				Data: map[string]interface{}{
					"email":   "",
					"id":      "some_uid",
					"name":    "",
					"phone":   "",
					"picture": "",
					"role":    float64(0),
				},
			},
		},
		{
			Name:           "Fail: User not found",
			ServiceReturns: nil,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusNotFound,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name:           "Fail: Unknown error",
			ServiceReturns: nil,
			ServiceErr:     errors.New("some error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "some error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockService.On("GetFullUserByID", mock.Anything, mock.Anything).Return(tc.ServiceReturns, tc.ServiceErr)

			s.fiberApp.Post("/:userID", func(ctx *fiber.Ctx) error {
				ctx.Locals("user", token) // set token to context
				return s.userController.GetFullUserByID(ctx)
			})
			r := httptest.NewRequest("POST", "/some_uid", nil)
			resp, err := s.fiberApp.Test(r)
			s.NoError(err)

			var body response.BaseResponse
			err = json.NewDecoder(resp.Body).Decode(&body)
			s.NoError(err)

			s.Equal(tc.ExpectedStatus, resp.StatusCode)
			s.Equal(tc.ExpectedBody, body)
		})
		s.TearDownTest()
	}
}
