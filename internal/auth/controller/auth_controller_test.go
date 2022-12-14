package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"office-booking-backend/internal/auth/dto"
	mockService "office-booking-backend/internal/auth/service/mock"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"testing"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestSuiteAuthController struct {
	suite.Suite
	mockAuthService *mockService.AuthServiceMock
	mockValidator   *validator.ValidatorMock
	authController  *AuthController
	fiberApp        *fiber.App
}

func (s *TestSuiteAuthController) SetupTest() {
	s.mockAuthService = new(mockService.AuthServiceMock)
	s.mockValidator = new(validator.ValidatorMock)
	s.authController = NewAuthController(s.mockAuthService, s.mockValidator)
	s.fiberApp = fiber.New(fiber.Config{
		ErrorHandler: response.DefaultErrorHandler,
	})
}

func (s *TestSuiteAuthController) TearDownTest() {
	s.mockAuthService = nil
	s.authController = nil
	s.fiberApp = nil
}

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestRegisterUser() {
	signupReq := dto.SignupRequest{
		Email:    "mail@mail.com",
		Password: "password",
		Name:     "name",
		Phone:    "phone",
	}

	for _, tc := range []struct {
		Name            string
		MimeType        string
		RequestBody     interface{}
		ValidatorReturn *validator.ErrorsResponse
		ServiceErr      error
		ExpectedStatus  int
		ExpectedBody    response.BaseResponse
		ExpectedErr     error
	}{
		{
			Name:           "Success registering user",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    signupReq,
			ExpectedStatus: fiber.StatusCreated,
			ExpectedBody: response.BaseResponse{
				Message: "user registered successfully",
				Data: map[string]interface{}{
					"uid": "",
				},
			},
		},
		{
			Name:           "Failed registering user: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("invalid request body"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name:        "Failed registering user: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: signupReq,
			ValidatorReturn: &validator.ErrorsResponse{{
				Field:  "email",
				Reason: "must be a valid email",
			}},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "email",
						"reason": "must be a valid email",
					},
				},
			},
		},
		{
			Name:           "Failed registering user: Email already exists",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    signupReq,
			ServiceErr:     err2.ErrDuplicateEmail,
			ExpectedStatus: fiber.StatusConflict,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrDuplicateEmail.Error(),
			},
		},
		{
			Name:           "Failed registering user: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    signupReq,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("RegisterUser", mock.Anything, mock.Anything).Return("", tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorReturn)

			s.fiberApp.Post("/register", s.authController.RegisterUser)
			r := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestLoginUser() {
	loginReq := dto.LoginRequest{
		Email:    "mail@mail.com",
		Password: "password",
	}

	token := &dto.TokenPair{
		AccessToken:  "some_access_token",
		RefreshToken: "some_refresh_token",
	}

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ValidatorError *validator.ErrorsResponse
		ServiceReturn  *dto.TokenPair
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
		ExpectedErr    error
	}{
		{
			Name:           "Success logging in user",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    loginReq,
			ServiceReturn:  token,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user logged in successfully",
				Data: map[string]interface{}{
					"role":         float64(0),
					"userId":       "",
					"accessExpAt":  float64(0),
					"accessToken":  token.AccessToken,
					"refreshExpAt": float64(0),
					"refreshToken": token.RefreshToken,
				},
			},
		},
		{
			Name:           "Failed logging in user: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("invalid request body"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name:        "Failed logging in user: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: loginReq,
			ValidatorError: &validator.ErrorsResponse{{
				Field:  "email",
				Reason: "must be a valid email",
			}},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "email",
						"reason": "must be a valid email",
					},
				},
			},
		},
		{
			Name:           "Failed logging in user: invalid credentials",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    loginReq,
			ServiceErr:     err2.ErrInvalidCredentials,
			ExpectedStatus: fiber.StatusUnauthorized,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidCredentials.Error(),
			},
		},
		{
			Name:           "Failed logging in user: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    loginReq,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("LoginUser", mock.Anything, mock.Anything).Return(tc.ServiceReturn, tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorError)

			s.fiberApp.Post("/", s.authController.LoginUser)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestLogoutUser() {
	// LogoutUser handler always get valid token from middleware
	token := &jwt.Token{
		Claims: jwt.MapClaims{
			"uid": "some_uid",
		},
	}

	for _, tc := range []struct {
		Name           string
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
		ExpectedErr    error
	}{
		{
			Name:           "Success logging out user",
			ServiceErr:     nil,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "user logged out successfully",
			},
		},
		{
			Name:           "Failed logging out user: service error",
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			s.mockAuthService.On("LogoutUser", mock.Anything, mock.Anything).Return(tc.ServiceErr)

			s.fiberApp.Post("/", func(ctx *fiber.Ctx) error {
				ctx.Locals("user", token) // set token to context
				return s.authController.LogoutUser(ctx)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestRefreshToken() {
	token := &dto.TokenPair{
		AccessToken:  "some_access_token",
		RefreshToken: "some_refresh_token",
	}

	reqBody := &dto.RefreshTokenRequest{
		RefreshToken: "some_refresh_token",
	}

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ValidatorError *validator.ErrorsResponse
		ServiceReturn  *dto.TokenPair
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
		ExpectedErr    error
	}{
		{
			Name:           "Success refreshing token",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    reqBody,
			ServiceReturn:  token,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "token refreshed successfully",
				Data: map[string]interface{}{
					"role":         float64(0),
					"userId":       "",
					"accessExpAt":  float64(0),
					"accessToken":  token.AccessToken,
					"refreshExpAt": float64(0),
					"refreshToken": token.RefreshToken,
				},
			},
		},
		{
			Name:           "Failed refreshing token: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("some invalid json"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: "invalid request body",
			},
		},
		{
			Name:        "Failed refreshing token: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: reqBody,
			ValidatorError: &validator.ErrorsResponse{
				{
					Field:  "refreshToken",
					Reason: "required",
				},
			},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "refreshToken",
						"reason": "required",
					},
				},
			},
		},
		{
			Name:           "Failed refreshing token: user not found",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    reqBody,
			ServiceReturn:  nil,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusUnauthorized,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name:           "Failed refreshing token: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    reqBody,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("RefreshToken", mock.Anything, mock.Anything).Return(tc.ServiceReturn, tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorError)

			s.fiberApp.Post("/", s.authController.RefreshToken)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestRequestOTP() {
	req := &dto.ResetPasswordOTPRequest{
		Email: "some_email",
	}

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ValidatorError *validator.ErrorsResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
		ExpectedErr    error
	}{
		{
			Name:           "Success requesting OTP",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "otp sent successfully",
			},
		},
		{
			Name:           "Failed requesting OTP: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("some invalid json"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name:        "Failed requesting OTP: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: req,
			ValidatorError: &validator.ErrorsResponse{
				{
					Field:  "email",
					Reason: "required",
				},
			},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "email",
						"reason": "required",
					},
				},
			},
		},
		{
			Name:           "Failed requesting OTP: user not found",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusNotFound,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name:           "Failed requesting OTP: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("RequestPasswordResetOTP", mock.Anything, mock.Anything).Return(tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorError)

			s.fiberApp.Post("/", s.authController.RequestPasswordResetOTP)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestVerifyOTP() {
	req := &dto.ResetPasswordOTPVerifyRequest{
		Email: "123@123.com",
		Code:  "123123",
	}

	key := "some_key"

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ValidatorError *validator.ErrorsResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
	}{
		{
			Name:           "Success verifying OTP",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "otp verified successfully",
				Data: map[string]interface{}{
					"key": key,
				},
			},
		},
		{
			Name:           "Failed verifying OTP: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("some invalid request"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name:        "Failed verifying OTP: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: req,
			ValidatorError: &validator.ErrorsResponse{
				{
					Field:  "email",
					Reason: "required",
				},
			},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "email",
						"reason": "required",
					},
				},
			},
		},
		{
			Name:           "Failed verifying OTP: Invalid OTP",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     err2.ErrInvalidOTP,
			ExpectedStatus: fiber.StatusUnauthorized,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidOTP.Error(),
			},
		},
		{
			Name:           "Failed verifying OTP: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("VerifyPasswordResetOTP", mock.Anything, mock.Anything).Return(&key, tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorError)

			s.fiberApp.Post("/", s.authController.VerifyPasswordResetOTP)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestResetPassword() {
	req := &dto.PasswordResetRequest{
		Email:    "123@123.com",
		Password: "123123123",
		Key:      "someKey",
	}

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ValidatorError *validator.ErrorsResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
	}{
		{
			Name:           "Success resetting password",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "password reset successfully",
			},
		},
		{
			Name:           "Failed resetting password: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("some invalid request"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name:        "Failed resetting password: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: req,
			ValidatorError: &validator.ErrorsResponse{
				{
					Field:  "email",
					Reason: "required",
				},
			},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "email",
						"reason": "required",
					},
				},
			},
		},
		{
			Name:           "Failed resetting password: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
		{
			Name:           "Failed resetting password: invalid key",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     err2.ErrInvalidOTPToken,
			ExpectedStatus: fiber.StatusUnauthorized,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidOTPToken.Error(),
			},
		},
		{
			Name:           "Failed resetting password: user not found",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusNotFound,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name:           "Failed resetting password: unable to update password",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("ResetPassword", mock.Anything, mock.Anything).Return(tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorError)

			s.fiberApp.Post("/", s.authController.ResetPassword)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

//goland:noinspection Annotator,Annotator
func (s *TestSuiteAuthController) TestChangePassword() {
	req := &dto.ChangePasswordRequest{
		OldPassword: "123123123",
		NewPassword: "456456456",
	}

	token := &jwt.Token{
		Claims: jwt.MapClaims{
			"uid": "123",
		},
	}

	for _, tc := range []struct {
		Name           string
		MimeType       string
		RequestBody    interface{}
		ValidatorError *validator.ErrorsResponse
		ServiceErr     error
		ExpectedStatus int
		ExpectedBody   response.BaseResponse
	}{
		{
			Name:           "Success changing password",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ExpectedStatus: fiber.StatusOK,
			ExpectedBody: response.BaseResponse{
				Message: "password changed successfully",
			},
		},
		{
			Name:           "Failed changing password: invalid request body",
			MimeType:       fiber.MIMETextPlain,
			RequestBody:    []byte("some invalid request"),
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
			},
		},
		{
			Name:        "Failed changing password: validation error",
			MimeType:    fiber.MIMEApplicationJSON,
			RequestBody: req,
			ValidatorError: &validator.ErrorsResponse{
				{
					Field:  "oldPassword",
					Reason: "required",
				},
			},
			ExpectedStatus: fiber.StatusBadRequest,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrInvalidRequestBody.Error(),
				Data: []interface{}{
					map[string]interface{}{
						"field":  "oldPassword",
						"reason": "required",
					},
				},
			},
		},
		{
			Name:           "Failed changing password: service error",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     errors.New("error"),
			ExpectedStatus: fiber.StatusInternalServerError,
			ExpectedBody: response.BaseResponse{
				Message: "error",
			},
		},
		{
			Name:           "Failed changing password: user not found",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     err2.ErrUserNotFound,
			ExpectedStatus: fiber.StatusNotFound,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrUserNotFound.Error(),
			},
		},
		{
			Name:           "Failed changing password: old password is incorrect",
			MimeType:       fiber.MIMEApplicationJSON,
			RequestBody:    req,
			ServiceErr:     err2.ErrPasswordNotMatch,
			ExpectedStatus: fiber.StatusConflict,
			ExpectedBody: response.BaseResponse{
				Message: err2.ErrPasswordNotMatch.Error(),
			},
		},
	} {
		s.SetupTest()
		s.Run(tc.Name, func() {
			jsonBody, err := json.Marshal(tc.RequestBody)
			s.NoError(err)

			s.mockAuthService.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(tc.ServiceErr)
			s.mockValidator.On("ValidateJSON", mock.Anything).Return(tc.ValidatorError)

			s.fiberApp.Put("/", func(ctx *fiber.Ctx) error {
				ctx.Locals("user", token)
				return s.authController.ChangePassword(ctx)
			})
			r := httptest.NewRequest("PUT", "/", bytes.NewBuffer(jsonBody))
			r.Header.Set(fiber.HeaderContentType, tc.MimeType)
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

func TestAuthController(t *testing.T) {
	suite.Run(t, new(TestSuiteAuthController))
}
