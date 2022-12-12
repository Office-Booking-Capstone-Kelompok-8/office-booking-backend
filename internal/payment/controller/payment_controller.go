package controller

import (
	"office-booking-backend/internal/payment/dto"
	"office-booking-backend/internal/payment/service"
	err2 "office-booking-backend/pkg/errors"
	"office-booking-backend/pkg/response"
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

func (p *PaymentController) GetBanks(c *fiber.Ctx) error {
	banks, err := p.service.GetBanks(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "banks retrieved successfully",
		Data:    banks,
	})
}

func (p *PaymentController) GetAllPaymentMethod(c *fiber.Ctx) error {
	payments, err := p.service.GetAllPaymentMethod(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "payments retrieved successfully",
		Data:    payments,
	})
}

func (p *PaymentController) GetPaymentMethodByID(c *fiber.Ctx) error {
	paymentID := c.Params("paymentMethodID")
	paymentIDInt, err := strconv.Atoi(paymentID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidPaymentID.Error())
	}

	payment, err := p.service.GetPaymentMethodByID(c.Context(), paymentIDInt)
	if err != nil {
		switch err {
		case err2.ErrPaymentNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "payment method retrieved successfully",
		Data:    payment,
	})
}

func (p *PaymentController) CreatePaymentMethod(c *fiber.Ctx) error {
	paymentRequest := new(dto.CreatePaymentRequest)
	err := c.BodyParser(paymentRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	errs := p.validator.ValidateJSON(*paymentRequest)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err = p.service.CreatePaymentMethod(c.Context(), paymentRequest)
	if err != nil {
		switch err {
		case err2.ErrInvalidBankID:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.BaseResponse{
		Message: "payment method created successfully",
	})
}

func (p *PaymentController) UpdatePaymentMethod(c *fiber.Ctx) error {
	paymentRequest := new(dto.UpdatePaymentRequest)
	err := c.BodyParser(paymentRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidRequestBody.Error())
	}

	paymentID := c.Params("paymentID")
	paymentIDInt, err := strconv.Atoi(paymentID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidPaymentID.Error())
	}

	errs := p.validator.ValidateJSON(*paymentRequest)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.BaseResponse{
			Message: err2.ErrInvalidRequestBody.Error(),
			Data:    errs,
		})
	}

	err = p.service.UpdatePaymentMethod(c.Context(), paymentIDInt, paymentRequest)
	if err != nil {
		switch err {
		case err2.ErrPaymentNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case err2.ErrInvalidBankID:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "payment method updated successfully",
	})
}

func (p *PaymentController) DeletePaymentMethod(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")
	paymentIDInt, err := strconv.Atoi(paymentID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err2.ErrInvalidPaymentID.Error())
	}

	err = p.service.DeletePaymentMethod(c.Context(), paymentIDInt)
	if err != nil {
		switch err {
		case err2.ErrPaymentNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(response.BaseResponse{
		Message: "payment method deleted successfully",
	})
}
