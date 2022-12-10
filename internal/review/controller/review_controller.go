package controller

import (
	"office-booking-backend/internal/review/dto"
	"office-booking-backend/internal/review/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/validator"

	"github.com/gofiber/fiber/v2"
)

type ReviewController struct {
	validator validator.Validator
	service   service.ReviewService
}

func NewReviewController(service service.ReviewService, validator validator.Validator) *ReviewController {
	return &ReviewController{
		validator: validator,
		service:   service,
	}
}

func (rv *ReviewController) CreateReservationReview(c *fiber.Ctx) error {
	reviewRequest := new(dto.AddReviewRequest)
	err := c.BodyParser(reviewRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := rv.validator.ValidateJSON(*reviewRequest)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err = rv.service.CreateReservationReview(c.Context(), reviewRequest)
	if err != nil {
		switch err {
		case err2.ErrInvalidBankID:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "Successfully created review",
	})
}
