package entity

import "gorm.io/gorm"

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
