package service

import (
	"context"
	"office-booking-backend/internal/payment/dto"
)

type PaymentService interface {
	GetPaymentByID(ctx context.Context, paymentID int) (*dto.PaymentResponse, error)
}
