package impl

import (
	"context"
	"log"
	"office-booking-backend/internal/payment/dto"
	"office-booking-backend/internal/payment/repository"
	"office-booking-backend/internal/payment/service"
)

type PaymentServiceImpl struct {
	repo repository.PaymentRepository
}

func NewPaymentServiceImpl(repo repository.PaymentRepository) service.PaymentService {
	return &PaymentServiceImpl{
		repo: repo,
	}
}

func (p PaymentServiceImpl) GetPaymentByID(ctx context.Context, paymentID int) (*dto.PaymentResponse, error) {
	payment, err := p.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		log.Println("error when get payment by id: ", err)
		return nil, err
	}

	return dto.NewPaymentResponse(payment), nil
}
