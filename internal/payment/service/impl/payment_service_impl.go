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

func (p *PaymentServiceImpl) GetBanks(ctx context.Context) (*dto.BanksResponse, error) {
	banks, err := p.repo.GetAllBank(ctx)
	if err != nil {
		log.Println("error when getting all banks: ", err)
		return nil, err
	}

	return dto.NewBanksResponse(banks), nil
}

func (p *PaymentServiceImpl) GetPaymentByID(ctx context.Context, paymentID int) (*dto.PaymentResponse, error) {
	payment, err := p.repo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		log.Println("error when get payment by id: ", err)
		return nil, err
	}

	return dto.NewPaymentResponse(payment), nil
}

func (p *PaymentServiceImpl) CreatePayment(ctx context.Context, payment *dto.CreatePaymentRequest) error {
	paymentEntity := payment.ToEntity()
	err := p.repo.CreatePayment(ctx, paymentEntity)
	if err != nil {
		log.Println("error when create payment: ", err)
		return err
	}

	return nil
}
func (p *PaymentServiceImpl) UpdatePayment(ctx context.Context, paymentID int, payment *dto.UpdatePaymentRequest) error {
	paymentEntity := payment.ToEntity()
	paymentEntity.ID = uint(paymentID)
	err := p.repo.UpdatePayment(ctx, paymentEntity)
	if err != nil {
		log.Println("error when update payment: ", err)
		return err
	}

	return nil
}
