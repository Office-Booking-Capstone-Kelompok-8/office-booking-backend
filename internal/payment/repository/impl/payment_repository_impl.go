package impl

import (
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"strings"

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
	// TODO: add implementation
	panic("implement me")
}

func (p PaymentRepositoryImpl) GetAllBank(ctx context.Context) (*entity.Banks, error) {
	banks := new(entity.Banks)
	err := p.db.WithContext(ctx).Find(banks).Error
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (p PaymentRepositoryImpl) GetPaymentByID(ctx context.Context, paymentID int) (*entity.Payment, error) {
	payment := new(entity.Payment)
	err := p.db.WithContext(ctx).
		Model(&entity.Payment{}).
		Joins("Bank").
		First(payment, paymentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err2.ErrPaymentNotFound
		}
		return nil, err
	}

	return payment, nil
}

func (p PaymentRepositoryImpl) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	err := p.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		if strings.Contains(err.Error(), "CONSTRAINT `fk_payments_bank`") {
			return err2.ErrInvalidBankID
		}

		return err
	}

	return nil
}

func (p PaymentRepositoryImpl) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	res := p.db.WithContext(ctx).Updates(payment)
	if res.Error != nil {
		if strings.Contains(res.Error.Error(), "CONSTRAINT `fk_payments_bank`") {
			return err2.ErrInvalidBankID
		}

		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrPaymentNotFound
	}

	return nil
}
