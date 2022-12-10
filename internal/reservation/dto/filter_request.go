package dto

import "office-booking-backend/pkg/custom"

type ReservationQueryParam struct {
	UserName     string      `query:"userName" validate:"omitempty,min=3,max=50"`
	UserID       string      `query:"userId" validate:"omitempty,uuid4"`
	BuildingID   string      `query:"buildingId" validate:"omitempty,uuid4"`
	StatusID     int         `query:"statusId" validate:"omitempty,gte=1,lte=6"`
	StartDate    custom.Date `query:"startDate"`
	EndDate      custom.Date `query:"endDate"`
	CreatedStart custom.Date `query:"createdStart" validate:"required_with=CreatedEnd"`
	CreatedEnd   custom.Date `query:"createdEnd" validate:"required_with=CreatedStart"`
	SortBy       string      `query:"sortBy" validate:"omitempty,oneof=created_at start_date end_date building_name user_name"`
	SortOrder    string      `query:"sortOrder" validate:"omitempty,oneof=asc desc"`
	Page         int         `query:"page" validate:"gte=1"`
	Limit        int         `query:"limit" validate:"gte=1"`
	Offset       int         `query:"-" validate:"isdefault"`
}
