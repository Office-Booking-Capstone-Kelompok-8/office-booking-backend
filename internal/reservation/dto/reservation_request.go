package dto

import (
	"office-booking-backend/pkg/entity"
	"time"
)

type AddReservartionRequest struct {
	BuildingID  string `json:"buildingId" validate:"required"`
	CompanyName string `json:"companyName" validate:"required"`
	StartDate   Date   `json:"startDate" validate:"required"`
	EndDate     Date   `json:"endDate" validate:"required,gtfield=StartDate"`
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
		EndDate:     a.EndDate.ToTime(),
	}
}
