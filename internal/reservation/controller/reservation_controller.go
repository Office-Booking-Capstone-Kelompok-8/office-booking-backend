package controller

import (
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"
	"reflect"
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
		Message: "reservations fetched successfully",
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
		Message: "reservations fetched successfully",
		Data:    reservations,
		Meta: fiber.Map{
			"total": count,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (r *ReservationController) GetReservationDetailByID(c *fiber.Ctx) error {
	reservationID := c.Params("ReservationID")
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
		Message: "reservation fetched successfully",
		Data:    reservation,
	})
}

func (r *ReservationController) GetUserReservationDetailByID(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("ReservationID")
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
		Message: "reservation fetched successfully",
		Data:    reservation,
	})
}

func (r *ReservationController) GetReservationTotal(c *fiber.Ctx) error {
	stat, err := r.service.GetReservationStat(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "reservation total fetched successfully",
		Data:    stat,
	})
}

func (r *ReservationController) GetTotalRevenueByTime(c *fiber.Ctx) error {
	stat, err := r.service.GetTotalRevenueByTime(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "reservation revenue fetched successfully",
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
		Message: "reservation created successfully",
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
		Message: "reservation created successfully",
		Data: fiber.Map{
			"reservationId": reservationID,
		},
	})
}

func (r *ReservationController) CancelReservation(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("ReservationID")

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
		Message: "reservation canceled successfully",
	})
}

func (r *ReservationController) UpdateReservation(c *fiber.Ctx) error {
	reservationID := c.Params("ReservationID")

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

	if reflect.DeepEqual(*reservation, dto.UpdateReservationRequest{}) {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
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
		Message: "reservation updated successfully",
	})
}

func (r *ReservationController) UpdateReservationStatus(c *fiber.Ctx) error {
	reservationID := c.Params("ReservationID")

	status := new(dto.UpdateReservationStatusRequest)
	if err := c.BodyParser(status); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateJSON(status); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err := r.service.UpdateReservationStatus(c.Context(), reservationID, status)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "reservation status updated successfully",
	})
}

func (r *ReservationController) DeleteReservation(c *fiber.Ctx) error {
	reservationID := c.Params("ReservationID")

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
		Message: "reservation deleted successfully",
	})
}

func (r *ReservationController) GetUserReservationReview(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("reservationID")

	review, err := r.service.GetReservationReview(c.Context(), reservationID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "reviews retrieved successfully",
		Data:    review,
	})
}

func (r *ReservationController) CreateReservationReview(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("reservationID")

	reviewRequest := new(dto.AddReviewRequest)
	if err := c.BodyParser(reviewRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateJSON(*reviewRequest); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err := r.service.CreateReservationReview(c.Context(), reviewRequest, reservationID, userID)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrNoPermission:
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		case err2.ErrReservationNotCompleted:
			fallthrough
		case err2.ErrReviewAlreadyExist:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "review created successfully",
	})
}

func (r *ReservationController) UpdateReservationReview(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["uid"].(string)

	reservationID := c.Params("reservationID")

	reviewRequest := new(dto.UpdateReviewRequest)
	if err := c.BodyParser(reviewRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	if errs := r.validator.ValidateJSON(*reviewRequest); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err := r.service.UpdateReservationReview(c.Context(), reviewRequest, reservationID, userID)
	if err != nil {
		switch err {
		case err2.ErrReservationNotFound:
			fallthrough
		case err2.ErrReviewNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrReviewNotEditable:
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "review updated successfully",
	})
}
