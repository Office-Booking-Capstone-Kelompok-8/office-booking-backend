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
