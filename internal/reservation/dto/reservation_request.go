package dto

import (
	"office-booking-backend/pkg/custom"
	"office-booking-backend/pkg/entity"
)

type AddAdminReservartionRequest struct {
	UserID      string      `json:"userId" validate:"required,uuid"`
	BuildingID  string      `json:"buildingId" validate:"required,uuid"`
	CompanyName string      `json:"companyName" validate:"required,min=3,max=255"`
	StartDate   custom.Date `json:"startDate" validate:"required"`
	Duration    int         `json:"duration" validate:"required,gte=1"`
}

type AddReservartionRequest struct {
	BuildingID  string      `json:"buildingId" validate:"required,uuid"`
	CompanyName string      `json:"companyName" validate:"required,min=3,max=255"`
	StartDate   custom.Date `json:"startDate" validate:"required"`
	Duration    int         `json:"duration" validate:"required,gte=1"`
}

func (a *AddReservartionRequest) ToEntity(userID string) *entity.Reservation {
	return &entity.Reservation{
		UserID:      userID,
		BuildingID:  a.BuildingID,
		CompanyName: a.CompanyName,
		StartDate:   a.StartDate.ToTime(),
		EndDate:     a.StartDate.ToTime().AddDate(0, a.Duration, 0),
	}
}

func (a *AddAdminReservartionRequest) ToEntity() *entity.Reservation {
	return &entity.Reservation{
		UserID:      a.UserID,
		BuildingID:  a.BuildingID,
		CompanyName: a.CompanyName,
		StartDate:   a.StartDate.ToTime(),
		EndDate:     a.StartDate.ToTime().AddDate(0, a.Duration, 0),
	}
}

type UpdateReservationRequest struct {
	UserID      string      `json:"userId" validate:"omitempty,uuid"`
	BuildingID  string      `json:"buildingId" validate:"omitempty,uuid"`
	CompanyName string      `json:"companyName" validate:"omitempty,min=3,max=255"`
	StartDate   custom.Date `json:"startDate" validate:"omitempty"`
	Duration    int         `json:"duration" validate:"required_with=StartDate"`
	StatusID    int         `json:"statusId" validate:"omitempty,gte=1,lte=6"`
	Message     string      `json:"message" validate:"omitempty,min=3,max=255"`
}

func (u *UpdateReservationRequest) ToEntity(reservationID string) *entity.Reservation {
	return &entity.Reservation{
		ID:          reservationID,
		UserID:      u.UserID,
		BuildingID:  u.BuildingID,
		CompanyName: u.CompanyName,
		StartDate:   u.StartDate.ToTime(),
		EndDate:     u.StartDate.ToTime().AddDate(0, u.Duration, 0),
		StatusID:    u.StatusID,
		Message:     u.Message,
	}
}

type AddReviewRequest struct {
	Rating  int    `json:"rating" validate:"required, gte=1, lte=5"`
	Message string `json:"comment" validate:"omitempty, min=3, max=255"`
}

func (a *AddReviewRequest) ToEntity() *entity.Review {
	return &entity.Review{
		Rating:  a.Rating,
		Message: a.Message,
	}
}
