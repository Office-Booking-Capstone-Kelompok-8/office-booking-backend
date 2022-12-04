package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	ID          string `gorm:"primaryKey; type:varchar(36); not null"`
	CompanyName string
	BuildingID  string `gorm:"type:varchar(36); not null" `
	Building    Building
	StartDate   time.Time `gorm:"type:datetime"`
	EndDate     time.Time `gorm:"type:datetime"`
	UserID      string    `gorm:"type:varchar(36); "`
	User        User      `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:SET NULL;"`
	StatusID    int       `gorm:"type:int; default:1"`
	Status      Status
}

func (r *Reservation) BeforeCreate(*gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}

type Reservations []Reservation

type Status struct {
	ID      int `gorm:"primaryKey; type:int; not null"`
	Message string
}
