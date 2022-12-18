package service

import (
	"context"
	"io"
	"office-booking-backend/internal/payment/dto"
)

type PaymentService interface {
	GetAllPaymentMethod(ctx context.Context) (*dto.PaymentMethodsResponse, error)
	GetPaymentMethodByID(ctx context.Context, paymentID int) (*dto.PaymentMethodResponse, error)
	GetBanks(ctx context.Context) (*dto.BanksResponse, error)
	CreatePaymentMethod(ctx context.Context, payment *dto.CreatePaymentRequest) (uint, error)
	AddPaymentProof(ctx context.Context, reservationID string, payment *dto.CreateReservationPaymentRequest, file io.Reader) error
	GetReservationPaymentByID(ctx context.Context, reservationID string) (*dto.PaymentDetailResponse, error)
	GetUserReservationPaymentByID(ctx context.Context, reservationID string, userID string) (*dto.PaymentDetailResponse, error)
	UpdatePaymentMethod(ctx context.Context, paymentID int, payment *dto.UpdatePaymentRequest) error
	DeletePaymentMethod(ctx context.Context, paymentID int) error
}
