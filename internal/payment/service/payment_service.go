package service

import (
	"context"
	"office-booking-backend/internal/payment/dto"
)

type PaymentService interface {
	GetAllPaymentMethod(ctx context.Context) (*dto.PaymentsResponse, error)
	GetPaymentMethodByID(ctx context.Context, paymentID int) (*dto.PaymentResponse, error)
	GetBanks(ctx context.Context) (*dto.BanksResponse, error)
	CreatePaymentMethod(ctx context.Context, payment *dto.CreatePaymentRequest) error
	UpdatePaymentMethod(ctx context.Context, paymentID int, payment *dto.UpdatePaymentRequest) error
	DeletePaymentMethod(ctx context.Context, paymentID int) error
}
