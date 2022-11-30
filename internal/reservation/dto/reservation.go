package dto

import (
	"office-booking-backend/pkg/entity"
	"time"
)

type BriefReservationResponse struct {
	ID        string `json:"id"`
	Building  *Building
	Tenant    *Tenant
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
}

func NewBriefReservationResponse(reservation *entity.Reservation) *BriefReservationResponse {
	return &BriefReservationResponse{
		ID: reservation.ID,
		Building: &Building{
			ID:   reservation.Building.ID,
			Name: reservation.Building.Name,
		},
		Tenant: &Tenant{
			ID:   reservation.User.ID,
			Name: reservation.User.Detail.Name,
		},
	}
}

type Building struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
