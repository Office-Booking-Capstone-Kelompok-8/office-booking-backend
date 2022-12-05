package dto

import (
	"office-booking-backend/pkg/entity"
	"time"
)

type AddAdminReservartionRequest struct {
	UserID      string `json:"userId" validate:"required,uuid"`
	BuildingID  string `json:"buildingId" validate:"required"`
	CompanyName string `json:"companyName" validate:"required"`
	StartDate   Date   `json:"startDate" validate:"required"`
	Duration    int    `json:"duration" validate:"required,gte=1"`
}

type AddReservartionRequest struct {
	BuildingID  string `json:"buildingId" validate:"required"`
	CompanyName string `json:"companyName" validate:"required"`
	StartDate   Date   `json:"startDate" validate:"required"`
	Duration    int    `json:"duration" validate:"required,gte=1"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	*d = Date{t}
	return nil
}

func (d *Date) ToTime() time.Time {
	return d.Time
}

func (a *AddReservartionRequest) ToEntity(userID string) *entity.Reservation {
	return &entity.Reservation{
		UserID:      userID,
		BuildingID:  a.BuildingID,
		CompanyName: a.CompanyName,
		StartDate:   a.StartDate.ToTime(),
		EndDate:     a.StartDate.AddDate(0, a.Duration, 0),
	}
}

func (a *AddAdminReservartionRequest) ToEntity() *entity.Reservation {
	return &entity.Reservation{
		UserID:      a.UserID,
		BuildingID:  a.BuildingID,
		CompanyName: a.CompanyName,
		StartDate:   a.StartDate.ToTime(),
		EndDate:     a.StartDate.AddDate(0, a.Duration, 0),
	}
}
