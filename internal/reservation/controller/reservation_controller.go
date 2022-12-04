package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"strconv"
)

type ReservationController struct {
	service   service.ReservationService
	validator validator.Validator
}

func NewReservationController(reservationService service.ReservationService, validator validator.Validator) *ReservationController {
	return &ReservationController{
		service:   reservationService,
		validator: validator,
	}
}

func (r *ReservationController) GetUserReservations(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	limit := c.Query("limit", "20")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	reservations, count, err := r.service.GetUserReservations(c.Context(), userID, pageInt, limitInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservations fetched successfully",
		Data:    reservations,
		Meta: fiber.Map{
			"total": count,
			"page":  pageInt,
			"limit": limitInt,
		},
	})
}

func (r *ReservationController) CreateReservation(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservation := new(dto.AddReservartionRequest)
	if err := c.BodyParser(reservation); err != nil {
		fmt.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateStruct(reservation); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	reservationID, err := r.service.CreateReservation(c.Context(), userID, reservation)
	if err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrBuildingNotAvailable:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "Reservation created successfully",
		Data: fiber.Map{
			"reservationID": reservationID,
		},
	})
}
