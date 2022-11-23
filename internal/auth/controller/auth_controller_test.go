package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"office-booking-backend/internal/auth/dto"
	mockService "office-booking-backend/internal/auth/service/mock"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/handler"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuiteAuthController struct {
	suite.Suite
	mockAuthService *mockService.AuthServiceMock
	authController  *AuthController
	fiberApp        *fiber.App
}

func (s *TestSuiteAuthController) SetupTest() {
	s.mockAuthService = new(mockService.AuthServiceMock)
	s.authController = NewAuthController(s.mockAuthService)
	s.fiberApp = fiber.New(fiber.Config{
		ErrorHandler: handler.DefaultErrorHandler,
	})
}

func (s *TestSuiteAuthController) TearDownTest() {
	s.mockAuthService = nil
	s.authController = nil
	s.fiberApp = nil
}

func (s *TestSuiteAuthController) TestRegisterUser() {
	signupReq := dto.SignupRequest{
		Email:    "mail@mail.com",
		Password: "password",
		Name:     "name",
		Phone:    "phone",
	}

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   fiber.Map
		ExpectedErr    error
	}{
		{
			Name:           "Success registering user",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    signupReq,
			ExpectedStatus: fiber.StatusCreated,
			ExpectedBody: fiber.Map{
				"message": "user created successfully",
			},
		},
		{
			Name:           "Failed registering user: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("invalid request body"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: fiber.Map{
				"message": "invalid request body",
			},
		},
		{
			Name:           "Failed registering user: Email already exists",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    signupReq,
			ServiceErr:     err2.ErrDuplicateEmail,
			ExpectedStatus: fiber.StatusConflict,
			ExpectedBody: fiber.Map{
				"message": err2.ErrDuplicateEmail.Error(),
			},
		},
		{
			Name:           "Failed registering user: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    signupReq,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: fiber.Map{
				"message": "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("RegisterUser", mock.Anything, mock.Anything).Return(tc.ServiceErr)

			s.fiberApp.Post("/register", s.authController.RegisterUser)
			r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
			resp, err := s.fiberApp.Test(r)
			s.NoError(err)

			var body fiber.Map
			err = json.NewDecoder(resp.Body).Decode(&body)
			s.NoError(err)

			s.Equal(tc.ExpectedStatus, resp.StatusCode)
			s.Equal(tc.ExpectedBody, body)
		})
		s.TearDownTest()
	}
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(TestSuiteAuthController))
}
