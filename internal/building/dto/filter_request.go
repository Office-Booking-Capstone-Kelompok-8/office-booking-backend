package dto

import "time"

type SearchBuildingQueryParam struct {
	BuildingName    string    `query:"buildingName" validate:"omitempty,min=3"`
	CityID          int       `query:"cityId" validate:"omitempty,uuid4"`
	DistrictID      int       `query:"districtId" validate:"omitempty,uuid4"`
	AnnualPriceMin  int       `query:"annualPriceMin" validate:"omitempty,gte=0"`
	AnnualPriceMax  int       `query:"annualPriceMax" validate:"omitempty,gte=0"`
	MonthlyPriceMin int       `query:"monthlyPriceMin" validate:"omitempty,gte=0"`
	MonthlyPriceMax int       `query:"monthlyPriceMax" validate:"omitempty,gte=0"`
	CapacityMin     int       `query:"capacityMin" validate:"omitempty,gte=0"`
	CapacityMax     int       `query:"capacityMax" validate:"omitempty,gte=0"`
	Latitude        float64   `query:"latitude" validate:"required_if=SortBy pinpoint"`
	Longitude       float64   `query:"longitude" validate:"required_if=SortBy pinpoint"`
	StartDate       time.Time `query:"-" validate:"required_with=Duration"`
	Duration        int       `query:"duration" validate:"required_with=StartDate"`
	EndDate         time.Time `query:"-"`
	SortBy          string    `query:"sortBy" validate:"omitempty,oneof=annual_price monthly_price capacity pinpoint"`
	Order           string    `query:"order" validate:"omitempty,oneof=asc desc"`
	Page            int       `query:"page" validate:"gte=1"`
	Limit           int       `query:"limit" validate:"gte=1"`
	Offset          int       `query:"-" validate:"isdefault"`
}
