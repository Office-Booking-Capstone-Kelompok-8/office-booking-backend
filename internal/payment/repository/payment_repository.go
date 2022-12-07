package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type PaymentRepository interface {
	GetAllPayment(ctx context.Context) (*entity.Payments, error)
	GetAllBank(ctx context.Context) (*entity.Banks, error)
	GetPaymentByID(ctx context.Context, paymentID int) (*entity.Payment, error)
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	UpdatePayment(ctx context.Context, payment *entity.Payment) error
}
