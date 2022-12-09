package dto

import "time"

type ReservationQueryParam struct {
	UserName   string    `query:"userName" validate:"omitempty,min=3,max=50"`
	UserID     string    `query:"userId" validate:"omitempty,uuid4"`
	BuildingID string    `query:"buildingId" validate:"omitempty,uuid4"`
	StatusID   int       `query:"statusId" validate:"omitempty,gte=1,lte=6"`
	StartDate  time.Time `query:"startDate"`
	EndDate    time.Time `query:"endDate"`
	Page       int       `query:"page" validate:"gte=1"`
	Limit      int       `query:"limit" validate:"gte=1"`
	Offset     int       `query:"-" validate:"isdefault"`
}
