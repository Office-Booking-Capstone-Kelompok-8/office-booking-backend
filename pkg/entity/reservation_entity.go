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
	Amount      int       `gorm:"type:int; not null"`
	UserID      string    `gorm:"type:varchar(36);"`
	User        User      `gorm:"constraint:OnUpdate:NO ACTION,OnDelete:SET NULL;"`
	StatusID    int       `gorm:"type:int; default:1"`
	Status      Status
	Message     string         `gorm:"type:varchar(255); default:''"`
	AcceptedAt  time.Time      `gorm:"type:datetime; default:NULL"`
	ExpiredAt   time.Time      `gorm:"type:datetime; default:NULL"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
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

// only used for returning stats
type StatusStat struct {
	StatusID   int64
	StatusName string
	Total      int64
}

type StatusesStat []StatusStat

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
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (rv *Review) BeforeCreate(*gorm.DB) (err error) {
	rv.ID = uuid.New().String()
	return
}

type Reviews []Review
