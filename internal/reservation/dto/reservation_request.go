package dto

import (
	"office-booking-backend/pkg/entity"
	"time"
)

type AddAdminReservartionRequest struct {
	UserID      string `json:"userId" validate:"required,uuid"`
	BuildingID  string `json:"buildingId" validate:"required,uuid"`
	CompanyName string `json:"companyName" validate:"required,min=3,max=255"`
	StartDate   Date   `json:"startDate" validate:"required"`
	Duration    int    `json:"duration" validate:"required,gte=1"`
}

type AddReservartionRequest struct {
	BuildingID  string `json:"buildingId" validate:"required,uuid"`
	CompanyName string `json:"companyName" validate:"required,min=3,max=255"`
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

type UpdateReservationRequest struct {
	UserID      string `json:"userId" validate:"omitempty,uuid"`
	BuildingID  string `json:"buildingId" validate:"omitempty,uuid"`
	CompanyName string `json:"companyName" validate:"omitempty,min=3,max=255"`
	StartDate   Date   `json:"startDate" validate:"omitempty"`
	Duration    int    `json:"duration" validate:"required_with=StartDate"`
	StatusID    int    `json:"statusId" validate:"omitempty,gte=1,lte=6"`
}

func (u *UpdateReservationRequest) ToEntity(reservationID string) *entity.Reservation {
	return &entity.Reservation{
		ID:          reservationID,
		UserID:      u.UserID,
		BuildingID:  u.BuildingID,
		CompanyName: u.CompanyName,
		StartDate:   u.StartDate.ToTime(),
		EndDate:     u.StartDate.AddDate(0, u.Duration, 0),
		StatusID:    u.StatusID,
	}
}
