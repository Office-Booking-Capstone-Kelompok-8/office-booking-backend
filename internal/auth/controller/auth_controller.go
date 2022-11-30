package controller

import (
	"errors"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthController struct {
	service   service.AuthService
	validator validator.Validator
}

func NewAuthController(AuthService service.AuthService, validator validator.Validator) *AuthController {
	return &AuthController{
		service:   AuthService,
		validator: validator,
	}
}

func (a *AuthController) RegisterUser(c *fiber.Ctx) error {
	user := new(dto.SignupRequest)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*user)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := a.service.RegisterUser(c.Context(), user); err != nil {
		if errors.Is(err, err2.ErrDuplicateEmail) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "user registered successfully",
	})
}

func (a *AuthController) RegisterAdmin(c *fiber.Ctx) error {
	user := new(dto.SignupRequest)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*user)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := a.service.RegisterAdmin(c.Context(), user); err != nil {
		if errors.Is(err, err2.ErrDuplicateEmail) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "admin registered successfully",
	})
}

func (a *AuthController) LoginUser(c *fiber.Ctx) error {
	user := new(dto.LoginRequest)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*user)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	tokenPair, err := a.service.LoginUser(c.Context(), user)
	if err != nil {
		if errors.Is(err, err2.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user logged in successfully",
		Data:    tokenPair,
	})
}

func (a *AuthController) LogoutUser(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	if err := a.service.LogoutUser(c.Context(), uid); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user logged out successfully",
	})
}

func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	token := new(dto.RefreshTokenRequest)
	if err := c.BodyParser(token); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*token)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	tokenPair, err := a.service.RefreshToken(c.Context(), token)
	if err != nil {
		if errors.Is(err, err2.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "token refreshed successfully",
		Data:    tokenPair,
	})
}

func (a *AuthController) RequestOTP(c *fiber.Ctx) error {
	otp := new(dto.OTPRequest)
	if err := c.BodyParser(otp); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*otp)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := a.service.RequestOTP(c.Context(), otp.Email); err != nil {
		if errors.Is(err, err2.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "otp sent successfully",
	})
}

func (a *AuthController) VerifyOTP(c *fiber.Ctx) error {
	otp := new(dto.OTPVerifyRequest)
	if err := c.BodyParser(otp); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*otp)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	key, err := a.service.VerifyOTP(c.Context(), otp)
	if err != nil {
		if errors.Is(err, err2.ErrInvalidOTP) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "otp verified successfully",
		Data: fiber.Map{
			"key": key,
		},
	})
}

func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	password := new(dto.PasswordResetRequest)
	if err := c.BodyParser(password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*password)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := a.service.ResetPassword(c.Context(), password); err != nil {
		if errors.Is(err, err2.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else if errors.Is(err, err2.ErrInvalidOTPToken) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "password reset successfully",
	})
}

func (a *AuthController) ChangePassword(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	password := new(dto.ChangePasswordRequest)
	if err := c.BodyParser(password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := a.validator.Validate(*password)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := a.service.ChangePassword(c.Context(), uid, password); err != nil {
		switch err {
		case err2.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrPasswordNotMatch:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "password changed successfully",
	})
}
