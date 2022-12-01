package dto

import (
	"office-booking-backend/pkg/entity"
)

type UpdateBuildingRequest struct {
	Name        string          `json:"name" validate:"omitempty,min=3,max=100"`
	Description string          `json:"description" validate:"omitempty,min=3,max=10000"`
	Capacity    int             `json:"capacity" validate:"omitempty,gte=1,lte=1000"`
	Prices      PriceRequest    `json:"price" validate:"omitempty,dive"`
	Owner       string          `json:"owner" validate:"omitempty"`
	Locations   LocationRequest `json:"location" validate:"omitempty,dive"`
}

func (c *UpdateBuildingRequest) ToEntity(buildingID string) *entity.Building {
	return &entity.Building{
		ID:           buildingID,
		Name:         c.Name,
		Description:  c.Description,
		Capacity:     c.Capacity,
		AnnualPrice:  c.Prices.AnnualPrice,
		MonthlyPrice: c.Prices.MonthlyPrice,
		Owner:        c.Owner,
		CityID:       c.Locations.CityID,
		DistrictID:   c.Locations.DistrictID,
	}
}

type PictureRequest struct {
	Index     int    `json:"index" validate:"required,gte=0,lte=10"`
	PictureID string `json:"pictureId" validate:"required"`
}

func (p *PictureRequest) ToEntity(buildingID string) *entity.Picture {
	return &entity.Picture{
		ID:         p.PictureID,
		BuildingID: buildingID,
		Index:      p.Index,
	}
}

type PicturesRequest []PictureRequest

func (p *PicturesRequest) ToEntity(buildingID string) *entity.Pictures {
	var pics entity.Pictures
	for _, picture := range *p {
		pics = append(pics, *picture.ToEntity(buildingID))
	}
	return &pics
}

type FacilityRequest struct {
	ID          int    `json:"id" validate:"omitempty"`
	Name        string `json:"name" validate:"required"`
	IconID      int    `json:"iconId" validate:"required"`
	Description string `json:"description" validate:"omitempty,min=3,max=1000"`
}

func (f *FacilityRequest) ToEntity(buildingID string) *entity.Facility {
	return &entity.Facility{
		ID:          f.ID,
		BuildingID:  buildingID,
		Name:        f.Name,
		CategoryID:  f.IconID,
		Description: f.Description,
	}
}

type FacilitiesRequest []FacilityRequest

func (f *FacilitiesRequest) ToEntity(buildingID string) *entity.Facilities {
	var facilities entity.Facilities
	for _, facility := range *f {
		facilities = append(facilities, *facility.ToEntity(buildingID))
	}
	return &facilities
}

type PriceRequest struct {
	AnnualPrice  int `json:"annual" validate:"omitempty,gte=0"`
	MonthlyPrice int `json:"monthly" validate:"omitempty,gte=0"`
}

type LocationRequest struct {
	Address    string `json:"address" validate:"omitempty,min=3,max=100"`
	DistrictID int    `json:"districtId" validate:"omitempty"`
	CityID     int    `json:"cityId" validate:"omitempty"`
	Geo        Geo    `json:"geo" validate:"omitempty,dive"`
}

type GeoRequest struct {
	Longitude float64 `json:"long" validate:"omitempty,latitude"`
	Latitude  float64 `json:"lat" validate:"omitempty,longitude"`
}
