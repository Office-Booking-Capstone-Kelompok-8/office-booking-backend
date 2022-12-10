package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID            string `gorm:"primaryKey; type:varchar(36); not null"`
	ReservationID string `gorm:"type:varchar(36); not null" `
	Reservation   Reservation
	UserID        string `gorm:"type:varchar(36);"`
	User          User   `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:SET NULL;"`
	BuildingID    string `gorm:"type:varchar(36); not null"`
	Building      Building
	Rating        int
	Message       string
	IsAnonymous   bool
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (r *Review) BeforeCreate(*gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}

type Reviews []Review
