package dto

import "office-booking-backend/pkg/entity"

type GetAllBuildingResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Pictures     entity.Pictures `json:"pictures"`
	AnnualPrice  int             `json:"annual_price"`
	MonthlyPrice int             `json:"monthly_price"`
	Owner        string          `json:"owner"`
	CityID       int             `json:"city_id"`
	DistrictID   int             `json:"district_id"`
}

func NewGetAllBuildingResponse(building *entity.Building) *GetAllBuildingResponse {
	return &GetAllBuildingResponse{
		ID:           building.ID,
		Name:         building.Name,
		Pictures:     building.Pictures,
		AnnualPrice:  building.AnnualPrice,
		MonthlyPrice: building.MonthlyPrice,
		Owner:        building.Owner,
		CityID:       building.CityID,
		DistrictID:   building.DistrictID,
	}
}

type GetAllBuildingsResponse []GetAllBuildingResponse

type Facility struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	IconName    string `json:"iconName"`
	Description string `json:"description"`
}

type Picture struct {
	ID  string `json:"id"`
	Url string `json:"url"`
	Alt string `json:"alt"`
}

type Price struct {
	AnnualPrice  int `json:"annual"`
	MonthlyPrice int `json:"monthly"`
}

type Location struct {
	Address    string `json:"address"`
	CityID     int    `json:"city"`
	DistrictID int    `json:"district"`
	Geos       Geo
}

type Geo struct {
	Longitude float64 `json:"long"`
	Latitude  float64 `json:"lat"`
}

type GetBuildingDetailByIDResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Pictures    Picture  `json:"pictures"`
	Description string   `json:"description"`
	Facilities  Facility `json:"facilities"`
	Capacity    int      `json:"capacity"`
	Prices      Price    `json:"price"`
	Owner       string   `json:"owner"`
	Locations   Location `json:"location"`
}

func NewGetBuildingDetailByIDResponse(building *entity.Building) *GetBuildingDetailByIDResponse {
	return &GetBuildingDetailByIDResponse{
		ID:         building.ID,
		Name:       building.Name,
		Pictures:   Picture{},
		Facilities: Facility{},
		Capacity:   building.Capacity,
		Prices:     Price{},
		Owner:      building.Owner,
		Locations:  Location{},
	}
}

type GetBuildingDetailsByIDResponse []GetBuildingDetailByIDResponse
