package impl

import (
	"fmt"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"strings"

	"github.com/Masterminds/squirrel"
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

func (p *PaymentRepositoryImpl) GetAllPaymentMethod(ctx context.Context) (*entity.Payments, error) {
	payment := new(entity.Payments)
	err := p.db.WithContext(ctx).
		Model(&entity.Payment{}).
		Joins("Bank").
		Find(payment).Error
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (p *PaymentRepositoryImpl) GetAllBank(ctx context.Context) (*entity.Banks, error) {
	banks := new(entity.Banks)
	err := p.db.WithContext(ctx).Find(banks).Error
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (p *PaymentRepositoryImpl) GetPaymentMethodByID(ctx context.Context, paymentID int) (*entity.Payment, error) {
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

func (p *PaymentRepositoryImpl) CreatePaymentMethod(ctx context.Context, payment *entity.Payment) error {
	err := p.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		if strings.Contains(err.Error(), "CONSTRAINT `fk_payments_bank`") {
			return err2.ErrInvalidBankID
		}

		return err
	}

	return nil
}

func (p *PaymentRepositoryImpl) CreateNewReservationPayment(ctx context.Context, payment *entity.Transaction) error {
	err := p.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		if strings.Contains(err.Error(), "CONSTRAINT `fk_transactions_payment`") {
			return err2.ErrPaymentMethodNotFound
		}
		return err
	}

	return nil
}

func (p *PaymentRepositoryImpl) GetReservationPaymentByID(ctx context.Context, reservationID string) (*entity.Transaction, error) {
	// rows, err := p.db.WithContext(ctx).
	// 	Table("transactions AS t").
	// 	Select("t.id, t.reservation_id, t.payment_id, t.proof_id, t.expired_at, t.paid_at, t.created_at, t.updated_at, p.id, p.account_name, p.account_number, p.account_name, b.icon, r.amount, r.start_date, r.end_date").
	// 	Joins("JOIN payments p ON p.id = t.payment_id").
	// 	Joins("JOIN banks b ON b.id = p.bank_id").
	// 	Joins("JOIN reservations r ON r.id = t.reservation_id").
	// 	Where("t.reservation_id = ?", reservationID).
	// 	Rows()
	db, err := p.db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := squirrel.Select("t.id, t.reservation_id, t.payment_id, t.proof_id, t.created_at, t.updated_at, p.id, p.account_name, p.account_number, p.account_name, b.icon, b.name, r.amount, r.start_date, r.end_date, pp.id, pp.url").
		From("transactions AS t").
		Join("reservations r ON r.id = t.reservation_id").
		Join("payments p ON p.id = t.payment_id").
		Join("banks b ON b.id = p.bank_id").
		Join("payment_proofs pp ON pp.id = t.proof_id").
		Where("t.reservation_id = ?", reservationID).
		RunWith(db).
		Query()

	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
	}()

	if rows.Next() {
		var tx entity.Transaction
		err = rows.Scan(&tx.ID, &tx.ReservationID, &tx.PaymentID, &tx.ProofID, &tx.CreatedAt, &tx.UpdatedAt, &tx.Payment.ID, &tx.Payment.AccountName, &tx.Payment.AccountNumber, &tx.Payment.AccountName, &tx.Payment.Bank.Icon, &tx.Payment.Bank.Name, &tx.Reservation.Amount, &tx.Reservation.StartDate, &tx.Reservation.EndDate, &tx.Proof.ID, &tx.Proof.URL)
		if err != nil {
			return nil, err
		}
		fmt.Println(tx)
		return &tx, nil
	}

	return nil, err2.ErrPaymentNotFound
}

func (p *PaymentRepositoryImpl) UpdatePaymentMethod(ctx context.Context, payment *entity.Payment) error {
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

func (p *PaymentRepositoryImpl) DeletePaymentMethod(ctx context.Context, paymentID int) error {
	res := p.db.WithContext(ctx).Delete(&entity.Payment{}, paymentID)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrPaymentNotFound
	}

	return nil
}
