package dto

import (
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/entity"
)

type BriefReservationResponse struct {
	ID          string                `json:"id"`
	Building    BriefBuildingResponse `json:"building"`
	CompanyName string                `json:"companyName"`
	StartDate   string                `json:"startDate"`
	EndDate     string                `json:"endDate"`
	Status      StatusResponse        `json:"status"`
}

func NewBriefReservationResponse(reservation *entity.Reservation) *BriefReservationResponse {
	return &BriefReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBriefBuildingResponse(reservation.Building),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Status:      *NewStatusResponse(reservation.Status),
	}
}

type BriefReservationsResponse []BriefReservationResponse

func NewBriefReservationsResponse(reservations *entity.Reservations) *BriefReservationsResponse {
	response := new(BriefReservationsResponse)
	for _, reservation := range *reservations {
		*response = append(*response, *NewBriefReservationResponse(&reservation))
	}
	return response
}

type BriefBuildingResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewBriefBuildingResponse(building entity.Building) *BriefBuildingResponse {
	return &BriefBuildingResponse{
		ID:   building.ID,
		Name: building.Name,
	}
}

type StatusResponse struct {
	ID      int    `json:"id"`
	Message string `json:"status"`
}

func NewStatusResponse(status entity.Status) *StatusResponse {
	return &StatusResponse{
		ID:      status.ID,
		Message: status.Message,
	}
}
