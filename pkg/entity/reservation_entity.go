package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	ID          string `gorm:"primaryKey; type:varchar(36); not null"`
	CompanyName string
	BuildingID  Building
	StartDate   time.Time `gorm:"type:datetime"`
	EndDate     time.Time `gorm:"type:datetime"`
	UserID      User
	StatusID    Status
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}

type Status struct {
	ID      int `gorm:"primaryKey; type:int; not null"`
	Message string
}
