package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type PaymentRepository interface {
	GetAllPaymentMethod(ctx context.Context) (*entity.Payments, error)
	GetAllBank(ctx context.Context) (*entity.Banks, error)
	GetPaymentMethodByID(ctx context.Context, paymentID int) (*entity.Payment, error)
	GetReservationPaymentByID(ctx context.Context, reservationID string) (*entity.Transaction, error)
	CreatePaymentMethod(ctx context.Context, payment *entity.Payment) error
	CreateReservationPayment(ctx context.Context, payment *entity.Transaction) error
	UpdatePaymentMethod(ctx context.Context, payment *entity.Payment) error
	DeletePaymentMethod(ctx context.Context, paymentID int) error
}
