package service

import (
	"context"
	"office-booking-backend/internal/payment/dto"
)

type PaymentService interface {
	GetPaymentByID(ctx context.Context, paymentID int) (*dto.PaymentResponse, error)
	GetBanks(ctx context.Context) (*dto.BanksResponse, error)
	CreatePayment(ctx context.Context, payment *dto.CreatePaymentRequest) error
	UpdatePayment(ctx context.Context, paymentID int, payment *dto.UpdatePaymentRequest) error
}
