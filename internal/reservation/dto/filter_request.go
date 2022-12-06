package dto

import "time"

type ReservationQueryParam struct {
	UserName   string `validate:"omitempty,min=3,max=50"`
	UserID     string `validate:"omitempty,uuid4"`
	BuildingID string `validate:"omitempty,uuid4"`
	StatusID   int    `validate:"omitempty,gte=1,lte=6"`
	StartDate  time.Time
	EndDate    time.Time
	Page       int `validate:"gte=1"`
	Offset     int
	Limit      int `validate:"gte=1"`
}
