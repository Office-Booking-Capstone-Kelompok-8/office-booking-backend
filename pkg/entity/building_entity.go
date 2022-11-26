package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Building struct {
	ID           string `gorm:"primaryKey; type:varchar(36); not null" `
	Name         string
	Description  string
	Capacity     int
	AnnualPrice  int
	MonthlyPrice int
	Owner        string
	Size         int
	CityID       City
	DistrictID   District
	Address      string
	Longitude    float64
	Latitude     float64
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (b *Building) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type District struct {
	ID     int `gorm:"primaryKey; type:int; not null"`
	CityID City
	Name   string
}

type City struct {
	ID   int `gorm:"primaryKey; type:int; not null"`
	Name string
}

type Picture struct {
	ID           string `gorm:"primaryKey; type:varchar(36); not null"`
	BuildingID   Building
	Url          string
	ThumbnailUrl string
	Alt          string
}

type Category struct {
	ID   int `gorm:"primaryKey; type:int; not null"`
	Name string
	Url  string
}

type Facility struct {
	ID          int `gorm:"primaryKey; type:int; not null"`
	BuildingID  Building
	CategoryID  Category
	Name        string
	Description string
}
