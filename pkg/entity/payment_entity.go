package entity

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	Name          string
	AccountNumber string
	Description   string
	BankID        string `gorm:"type:varchar(36)"`
	Bank          Bank   `gorm:"foreignKey:BankID"`
}

type Payments []Payment

type Bank struct {
	ID   int `gorm:"primaryKey; not null"`
	Name string
	Icon string
}
