package controller

import (
	"errors"
	"office-booking-backend/internal/auth/dto"
	"office-booking-backend/internal/auth/service"
	err2 "office-booking-backend/pkg/errors"

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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
	})
}
