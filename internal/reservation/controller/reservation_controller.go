package controller

import (
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

func (r *ReservationController) GetReservations(c *fiber.Ctx) error {
	filter := &dto.ReservationQueryParam{}
	if err := c.QueryParser(filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidQueryParams.Error())
	}

	if errs := r.validator.ValidateQuery(filter); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidQueryParams.Error(),
			Data:    errs,
		})
	}

	reservations, count, err := r.service.GetReservations(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservations fetched successfully",
		Data:    reservations,
		Meta: fiber.Map{
			"total": count,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (r *ReservationController) GetReservationDetailByID(c *fiber.Ctx) error {
	reservationID := c.Params("reservationID")
	reservation, err := r.service.GetReservationByID(c.Context(), reservationID)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservation fetched successfully",
		Data:    reservation,
	})
}

func (r *ReservationController) GetUserReservationDetailByID(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("reservationID")
	reservation, err := r.service.GetUserReservationByID(c.Context(), reservationID, userID)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservation fetched successfully",
		Data:    reservation,
	})
}

func (r *ReservationController) GetReservationTotal(c *fiber.Ctx) error {
	stat, err := r.service.GetReservationStat(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservation total fetched successfully",
		Data:    stat,
	})
}

func (r *ReservationController) CreateReservation(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservation := new(dto.AddReservartionRequest)
	if err := c.BodyParser(reservation); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateJSON(reservation); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	reservationID, err := r.service.CreateReservation(c.Context(), userID, reservation)
	if err != nil {
		switch err {
		case err2.ErrStartDateBeforeToday:
			fallthrough
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidBuildingID.Error())
		case err2.ErrBuildingNotAvailable:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		case err2.ErrInvalidUserID:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "Reservation created successfully",
		Data: fiber.Map{
			"reservationId": reservationID,
		},
	})
}

func (r *ReservationController) CreateAdminReservation(c *fiber.Ctx) error {
	reservation := new(dto.AddAdminReservartionRequest)
	if err := c.BodyParser(reservation); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateJSON(reservation); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	reservationID, err := r.service.CreateAdminReservation(c.Context(), reservation)
	if err != nil {
		switch err {
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidBuildingID.Error())
		case err2.ErrBuildingNotAvailable:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "Reservation created successfully",
		Data: fiber.Map{
			"reservationId": reservationID,
		},
	})
}

func (r *ReservationController) CancelReservation(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("reservationID")

	err := r.service.CancelReservation(c.Context(), userID, reservationID)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrNoPermission:
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		case err2.ErrReservationActive:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservation canceled successfully",
	})
}

func (r *ReservationController) UpdateReservation(c *fiber.Ctx) error {
	reservationID := c.Params("reservationID")

	reservation := new(dto.UpdateReservationRequest)
	if err := c.BodyParser(reservation); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateJSON(reservation); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err := r.service.UpdateReservation(c.Context(), reservationID, reservation)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrBuildingNotAvailable:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		case err2.ErrInvalidStatus:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case err2.ErrBuildingNotFound:
			return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidBuildingID.Error())
		case err2.ErrUserNotFound:
			return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidUserID.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reservation updated successfully",
	})
}

func (r *ReservationController) DeleteReservation(c *fiber.Ctx) error {
	reservationID := c.Params("reservationID")

	err := r.service.DeleteReservationByID(c.Context(), reservationID)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Review deleted successfully",
	})
}

func (r *ReservationController) GetReservationReviews(c *fiber.Ctx) error {
	reviews, err := r.service.GetReservationReviews(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "Reviews retrieved successfully",
		Data:    reviews,
	})
}
