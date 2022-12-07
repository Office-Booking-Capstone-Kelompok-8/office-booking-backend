package controller

import (
	"office-booking-backend/internal/payment/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	validator validator.Validator
	service   service.PaymentService
}

func NewPaymentController(service service.PaymentService, validator validator.Validator) *PaymentController {
	return &PaymentController{
		validator: validator,
		service:   service,
	}
}

func (p *PaymentController) GetPaymentByID(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")
	paymentIDInt, err := strconv.Atoi(paymentID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidPaymentID.Error())
	}

	payment, err := p.service.GetPaymentByID(c.Context(), paymentIDInt)
	if err != nil {
		switch err {
		case err2.ErrPaymentNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success get payment by id",
		"data":    payment,
	})
}
