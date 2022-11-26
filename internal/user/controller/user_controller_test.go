package controller

import (
	"bytes"
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

func (s *TestSuiteUserController) TestGetAllUsers() {
	for _, tc := range []struct {
		Name           string
		Q              string
		Limit          string
		Page           string
		ServiceReturns *dto.BriefUsersResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
	}{
		{
			Name: "success",
			ServiceReturns: &dto.BriefUsersResponse{
				{
					ID: "some_uid",
				},
			},
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user fetched successfully",
				Data: []interface{}{
					map[string]interface{}{
						"email": "",
						"id":    "some_uid",
						"name":  "",
						"phone": "",
					},
				},
				Meta: map[string]interface{}{
					"limit": float64(20),
					"page":  float64(1),
					"total": float64(1),
				},
			},
		},
		{
			Name:           "Success: No user",
			ServiceReturns: &dto.BriefUsersResponse{},
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user fetched successfully",
				Data:    []interface{}{},
				Meta: map[string]interface{}{
					"limit": float64(20),
					"page":  float64(1),
					"total": float64(1),
				},
			},
		},
		{
			Name:           "Fail: Limit is not a number",
			Limit:          "some_limit",
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidQueryParams.Error(),
			},
		},
		{
			Name:           "Fail: Page is not a number",
			Page:           "some_page",
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidQueryParams.Error(),
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
			s.mockService.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.ServiceReturns, int64(1), tc.ServiceErr)

			s.fiberApp.Get("/", func(ctx *fiber.Ctx) error {
				ctx.Context().QueryArgs().Set("q", tc.Q)
				ctx.Context().QueryArgs().Set("limit", tc.Limit)
				ctx.Context().QueryArgs().Set("page", tc.Page)
				return s.userController.GetAllUsers(ctx)
			})
			r := httptest.NewRequest("GET", "/", nil)
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

func (s *TestSuiteUserController) TestUpdateUser() {
	for _, tc := range []struct {
		Name           string
		Request        interface{}
		Mime           string
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
	}{
		{
			Name: "success",
			Request: dto.UserUpdateRequest{
				Name:  "some_name",
				Phone: "some_phone",
			},
			Mime:           fiber.MIMEApplicationJSON,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user updated successfully",
			},
		},
		{
			Name: "Fail: Invalid request",
			Request: dto.UserUpdateRequest{
				Name:  "some_name",
				Phone: "some_phone",
			},
			Mime:           fiber.MIMEApplicationXML,
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name: "Fail: user not found",
			Request: dto.UserUpdateRequest{
				Name:  "some_name",
				Phone: "some_phone",
			},
			Mime:           fiber.MIMEApplicationJSON,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusNotFound,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name: "Fail: Unknown error",
			Request: dto.UserUpdateRequest{
				Name:  "some_name",
				Phone: "some_phone",
			},
			Mime:           fiber.MIMEApplicationJSON,
			ServiceErr:     errors.New("some error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "some error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockService.On("UpdateUserByID", mock.Anything, mock.Anything, mock.Anything).Return(tc.ServiceErr)

			s.fiberApp.Put("/:userID", s.userController.UpdateUserByID)

			jsonBody := new(bytes.Buffer)
			err := json.NewEncoder(jsonBody).Encode(tc.Request)
			s.NoError(err)

			r := httptest.NewRequest("PUT", "/some_uid", jsonBody)
			r.Header.Set("Content-Type", tc.Mime)
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
