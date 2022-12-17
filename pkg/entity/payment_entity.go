package entity

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	AccountName   string
	AccountNumber string
	Description   string
	BankID        int
	Bank          Bank `gorm:"foreignKey:BankID"`
}

type Payments []Payment

type Bank struct {
	ID   int `gorm:"primaryKey; not null"`
	Name string
	Icon string
}

type Banks []Bank

type Transaction struct {
	ID            string `gorm:"primaryKey; type:varchar(36); not null"`
	ReservationID string `gorm:"type:varchar(36); not null"`
	Reservation   Reservation
	PaymentID     uint
	Payment       Payment
	ProofID       string
	Proof         PaymentProof
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (t *Transaction) BeforeCreate(*gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}

type PaymentProof struct {
	ID        string `gorm:"primaryKey; type:varchar(36); not null"`
	URL       string
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeleteAt  gorm.DeletedAt `gorm:"index"`
}
