package controller

import (
	"errors"
	"office-booking-backend/internal/user/dto"
	"office-booking-backend/internal/user/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserController struct {
	userService service.UserService
	validator   validator.Validator
}

func NewUserController(userService service.UserService, validator validator.Validator) *UserController {
	return &UserController{
		userService: userService,
		validator:   validator,
	}
}

func (u *UserController) GetLoggedFullUserByID(c *fiber.Ctx) error {
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

func (u *UserController) GetFullUserByID(c *fiber.Ctx) error {
	uid := c.Params("userID")

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

func (u *UserController) GetAllUsers(c *fiber.Ctx) error {
	q := c.Query("q")
	limit := c.Query("limit", "20")
	page := c.Query("page", "1")
	role := c.Query("role") // 1 - user, 2 - admin

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	roleInt, err := strconv.Atoi(role)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	users, total, err := u.userService.GetAllUsers(c.Context(), q, roleInt, limitInt, pageInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user fetched successfully",
		Data:    users,
		Meta: fiber.Map{
			"limit": limitInt,
			"page":  pageInt,
			"total": total,
		},
	})
}

func (u *UserController) UpdateUserByID(c *fiber.Ctx) error {
	uid := c.Params("userID")

	user := new(dto.UserUpdateRequest)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := u.validator.ValidateStruct(*user); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := u.userService.UpdateUserByID(c.Context(), uid, user); err != nil {
		switch err {
		case err2.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrDuplicateEmail:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user updated successfully",
	})
}

func (u *UserController) UpdateLoggedUser(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	user := new(dto.UserUpdateRequest)
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := u.validator.ValidateStruct(*user); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	if err := u.userService.UpdateUserByID(c.Context(), uid, user); err != nil {
		switch err {
		case err2.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrDuplicateEmail:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user updated successfully",
	})
}

func (u *UserController) DeleteUserByID(c *fiber.Ctx) error {
	uid := c.Params("userID")

	if err := u.userService.DeleteUserByID(c.Context(), uid); err != nil {
		switch err {
		case err2.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrUserHasReservation:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user deleted successfully",
	})
}

func (u *UserController) UpdateUserAvatar(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	file, err := form.File["picture"][0].Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if err := u.userService.UploadUserAvatar(c.Context(), uid, file); err != nil {
		switch err {
		case err2.ErrUserNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	err = file.Close()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user avatar updated successfully",
	})
}
