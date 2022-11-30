package controller

import (
	"office-booking-backend/internal/reservation/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReservationsController struct {
	reservationsService service.ReservationsService
	validator           validator.Validator
}

func NewReservationController(reservationsService service.ReservationsService, validator validator.Validator) *ReservationsController {
	return &ReservationsController{
		reservationsService: reservationsService,
		validator:           validator,
	}
}

func (r *ReservationsController) GetAllReservations(c *fiber.Ctx) error {
	status := c.Query("status")
	buildingID := c.Query("buildingID")
	userID := c.Query("userID")
	userName := c.Query("userName")
	startDate := c.Query("startDate", "0001-01-01")
	endDate := c.Query("endDate", "0001-01-01")
	page := c.Query("page", "1")
	limit := c.Query("limit", "20")

	// Parse startDate to time.Time YYYY-MM-DD format
	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	// Parse endDate to time.Time YYYY-MM-DD format
	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	reservations, total, err := r.reservationsService.GetAllReservations(c.Context(), status, buildingID, userID, userName, startDateParsed, endDateParsed, limitInt, pageInt)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "user fetched successfully",
		Data:    reservations,
		Meta: fiber.Map{
			"limit": limitInt,
			"page":  pageInt,
			"total": total,
		},
	})
}
