package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Building struct {
	ID           string `gorm:"primaryKey; type:varchar(36); not null" `
	Name         string
	Description  string   `gorm:"type:text"`
	Pictures     Pictures `gorm:"foreignKey:BuildingID"`
	Capacity     int
	AnnualPrice  int
	MonthlyPrice int
	Facilities   Facilities `gorm:"foreignKey:BuildingID"`
	Owner        string
	Size         int
	CityID       int `gorm:"default:null"`
	City         City
	DistrictID   int `gorm:"default:null"`
	District     District
	Address      string `gorm:"type:text"`
	Longitude    float64
	Latitude     float64
	CreatedByID  string         `gorm:"type:varchar(36); default:null;"`
	CreatedBy    User           `gorm:"foreignKey:CreatedByID"`
	IsPublished  *bool          `gorm:"default:false"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Buildings []Building

func (b *Building) BeforeCreate(*gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

type District struct {
	ID     int `gorm:"primaryKey; type:int; not null"`
	CityID int
	Name   string
}

type Districts []District

type City struct {
	ID        int `gorm:"primaryKey; type:int; not null"`
	Districts Districts
	Name      string
}

type Cities []City

type Picture struct {
	ID           string `gorm:"primaryKey; type:varchar(36); not null"`
	Key          string `gorm:"type:varchar(36); not null"`
	BuildingID   string `gorm:"type:varchar(36); not null" `
	Index        *int
	Url          string
	ThumbnailUrl string
	Alt          string
}

type Pictures []Picture

type Category struct {
	ID   int `gorm:"primaryKey; type:int; not null"`
	Name string
	Url  string
}

type Categories []Category

type Facility struct {
	ID          int    `gorm:"primaryKey; type:int; not null"`
	BuildingID  string `gorm:"type:varchar(36); not null" `
	CategoryID  int    `gorm:"type:int; not null"`
	Category    Category
	Name        string
	Description string
}

type Facilities []Facility
