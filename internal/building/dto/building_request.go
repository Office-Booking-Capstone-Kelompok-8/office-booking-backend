package dto

import (
	"office-booking-backend/pkg/entity"
)

type AddFacilityRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	IconID      int    `json:"iconId" validate:"required"`
	Description string `json:"description" validate:"omitempty,min=3,max=1000"`
}

func (f *AddFacilityRequest) ToEntity(buildingID string) *entity.Facility {
	return &entity.Facility{
		BuildingID:  buildingID,
		Name:        f.Name,
		CategoryID:  f.IconID,
		Description: f.Description,
	}
}

type AddFacilitiesRequest []AddFacilityRequest

func (f *AddFacilitiesRequest) ToEntity(buildingID string) *entity.Facilities {
	var facilities entity.Facilities
	for _, facility := range *f {
		facilities = append(facilities, *facility.ToEntity(buildingID))
	}
	return &facilities
}

type UpdateBuildingRequest struct {
	Name        string                  `json:"name" validate:"omitempty,min=3,max=100"`
	Description string                  `json:"description" validate:"omitempty,min=3,max=10000"`
	Facilities  UpdateFacilitiesRequest `json:"facilities" validate:"omitempty,dive"`
	Pictures    PicturesRequest         `json:"pictures" validate:"omitempty,dive"`
	Capacity    int                     `json:"capacity" validate:"omitempty,gte=1"`
	Size        int                     `json:"size" validate:"omitempty,gte=1"`
	Prices      PriceRequest            `json:"price" validate:"omitempty,dive"`
	Owner       string                  `json:"owner" validate:"omitempty"`
	Locations   LocationRequest         `json:"location" validate:"omitempty,dive"`
	IsPublished bool                    `json:"isPublished" validate:"omitempty"`
}

func (c *UpdateBuildingRequest) ToEntity(buildingID string) *entity.Building {
	return &entity.Building{
		ID:           buildingID,
		Name:         c.Name,
		Description:  c.Description,
		Pictures:     *c.Pictures.ToEntity(buildingID),
		Capacity:     c.Capacity,
		AnnualPrice:  c.Prices.AnnualPrice,
		MonthlyPrice: c.Prices.MonthlyPrice,
		Facilities:   *c.Facilities.ToEntity(buildingID),
		Owner:        c.Owner,
		Size:         c.Size,
		CityID:       c.Locations.CityID,
		DistrictID:   c.Locations.DistrictID,
		Address:      c.Locations.Address,
		Longitude:    c.Locations.Geo.Longitude,
		Latitude:     c.Locations.Geo.Latitude,
		IsPublished:  &c.IsPublished,
	}
}

type PictureRequest struct {
	Index     int    `json:"index" validate:"gte=0,lte=9"`
	PictureID string `json:"pictureId" validate:"required"`
}

func (p *PictureRequest) ToEntity(buildingID string) *entity.Picture {
	return &entity.Picture{
		ID:         p.PictureID,
		BuildingID: buildingID,
		Index:      &p.Index,
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

type UpdateFacilityRequest struct {
	ID          int    `json:"facilityId" validate:"required"`
	Name        string `json:"name" validate:"omitempty,min=3,max=100"`
	IconID      int    `json:"iconId" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty,min=3,max=1000"`
}

func (f *UpdateFacilityRequest) ToEntity(buildingID string) *entity.Facility {
	return &entity.Facility{
		ID:          f.ID,
		BuildingID:  buildingID,
		Name:        f.Name,
		CategoryID:  f.IconID,
		Description: f.Description,
	}
}

type UpdateFacilitiesRequest []UpdateFacilityRequest

func (f *UpdateFacilitiesRequest) ToEntity(buildingID string) *entity.Facilities {
	var facilities entity.Facilities
	for _, facility := range *f {
		facilities = append(facilities, *facility.ToEntity(buildingID))
	}
	return &facilities
}

type PriceRequest struct {
	AnnualPrice  int `json:"annual" validate:"omitempty,gte=1"`
	MonthlyPrice int `json:"monthly" validate:"omitempty,gte=1"`
}

type LocationRequest struct {
	Address    string `json:"address" validate:"omitempty,min=3"`
	DistrictID int    `json:"districtId" validate:"omitempty"`
	CityID     int    `json:"cityId" validate:"omitempty"`
	Geo        Geo    `json:"geo" validate:"omitempty,dive"`
}

type GeoRequest struct {
	Longitude float64 `json:"long" validate:"omitempty,latitude"`
	Latitude  float64 `json:"lat" validate:"omitempty,longitude"`
}
