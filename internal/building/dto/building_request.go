package dto

import "office-booking-backend/pkg/entity"

type CreateBuildingRequest struct {
	ID          string            `json:"id" validate:"required,uuid"`
	Name        string            `json:"name" validate:"required,min=3,max=100"`
	Pictures    PicturesRequest   `json:"pictures" validate:"required"`
	Description string            `json:"description" validate:"required,min=3,max=1000"`
	Facilities  FacilitiesRequest `json:"facilities" validate:"required"`
	Capacity    int               `json:"capacity" validate:"required,gte=1,lte=1000"`
	Prices      PriceRequest      `json:"price" validate:"required"`
	Owner       string            `json:"owner" validate:"required"`
	Locations   LocationRequest   `json:"location" validate:"required"`
}

func (c *CreateBuildingRequest) ToEntity() *entity.Building {
	return &entity.Building{
		ID:           c.ID,
		Name:         c.Name,
		Pictures:     *c.Pictures.ToEntity(),
		Description:  c.Description,
		Facilities:   *c.Facilities.ToEntity(),
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

func (p *PictureRequest) ToEntity() *entity.Picture {
	return &entity.Picture{
		ID:    p.PictureID,
		Index: p.Index,
	}
}

type PicturesRequest []PictureRequest

func (p *PicturesRequest) ToEntity() *entity.Pictures {
	var pics entity.Pictures
	for _, picture := range *p {
		pics = append(pics, *picture.ToEntity())
	}
	return &pics
}

type FacilityRequest struct {
	Name        string `json:"name" validate:"required"`
	IconID      int    `json:"iconId" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (f *FacilityRequest) ToEntity() *entity.Facility {
	return &entity.Facility{
		Name:        f.Name,
		CategoryID:  f.IconID,
		Description: f.Description,
	}
}

type FacilitiesRequest []FacilityRequest

func (f *FacilitiesRequest) ToEntity() *entity.Facilities {
	var facilities entity.Facilities
	for _, facility := range *f {
		facilities = append(facilities, *facility.ToEntity())
	}
	return &facilities
}

type PriceRequest struct {
	AnnualPrice  int `json:"annual" validate:"required,gte=0"`
	MonthlyPrice int `json:"monthly" validate:"required,gte=0"`
}

type LocationRequest struct {
	Address    string `json:"address" validate:"required,min=3,max=100"`
	DistrictID int    `json:"districtId" validate:"required"`
	CityID     int    `json:"cityId" validate:"required"`
	Geo        Geo    `json:"geo" validate:"required"`
}

type GeoRequest struct {
	Longitude float64 `json:"long" validate:"required,latitude"`
	Latitude  float64 `json:"lat" validate:"required,longitude"`
}
