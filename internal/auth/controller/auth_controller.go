package controller

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(AuthService service.AuthService) *AuthController {
	return &AuthController{
		service: AuthService,
	}
}

func (a *AuthController) RegisterUser(c *fiber.Ctx) error {
	var user dto.SignupRequest
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if err := a.service.RegisterUser(c.Context(), &user); err != nil {
		if errors.Is(err, err2.ErrDuplicateEmail) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "user registered successfully",
	})
}

func (a *AuthController) LoginUser(c *fiber.Ctx) error {
	var user dto.LoginRequest
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	tokenPair, err := a.service.LoginUser(c.Context(), &user)
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
