package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/internal/user/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) GetFullUserByID(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	user, err := u.userService.GetFullUserByID(c.Context(), uid)
	if err != nil {
		if errors.Is(err, err2.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user fetched successfully",
		Data:    user,
	})
}
