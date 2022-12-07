package impl

import (
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepositoryImpl(db *gorm.DB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{
		db: db,
	}
}

func (p PaymentRepositoryImpl) GetAllPayment(ctx context.Context) (*entity.Payments, error) {
	//TODO implement me
	panic("implement me")
}

func (p PaymentRepositoryImpl) GetPaymentByID(ctx context.Context, paymentID int) (*entity.Payment, error) {
	payment := new(entity.Payment)
	err := p.db.WithContext(ctx).First(payment, paymentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrPaymentNotFound
		}
		return nil, err
	}

	return payment, nil
}

func (p PaymentRepositoryImpl) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (p PaymentRepositoryImpl) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (p PaymentRepositoryImpl) DeletePayment(ctx context.Context, paymentID int) error {
	//TODO implement me
	panic("implement me")
}
