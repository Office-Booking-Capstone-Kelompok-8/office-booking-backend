package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID            int `gorm:"primaryKey; type:int; not null"`
	Name          string
	AccountNumber string
	Description   string
	LogoID        PaymentPicture
}

type PaymentPicture struct {
	ID  string `gorm:"primaryKey; type:varchar(36); not null"`
	Url int
	Alt string
}

func (pp *PaymentPicture) BeforeCreate(tx *gorm.DB) (err error) {
	pp.ID = uuid.New().String()
	return
}
