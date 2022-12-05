package dto

import (
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/entity"
	"time"
)

type BriefReservationResponse struct {
	ID          string                `json:"id"`
	Building    BriefBuildingResponse `json:"building"`
	CompanyName string                `json:"companyName"`
	StartDate   string                `json:"startDate"`
	EndDate     string                `json:"endDate"`
	Status      StatusResponse        `json:"status"`
	CreatedAt   time.Time             `json:"createdAt"`
}

func NewBriefReservationResponse(reservation *entity.Reservation) *BriefReservationResponse {
	return &BriefReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBriefBuildingResponse(reservation.Building),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Status:      *NewStatusResponse(reservation.Status),
		CreatedAt:   reservation.CreatedAt,
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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func NewBriefBuildingResponse(building entity.Building) *BriefBuildingResponse {
	return &BriefBuildingResponse{
		ID:       building.ID,
		Name:     building.Name,
		Picture:  building.Pictures[0].ThumbnailUrl,
		City:     building.City.Name,
		District: building.District.Name,
		Address:  building.Address,
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

type FullAdminReservationResponse struct {
	ID          string                `json:"id"`
	Building    BriefBuildingResponse `json:"building"`
	Tenant      BriefTenantResponse   `json:"tenant"`
	CompanyName string                `json:"companyName"`
	StartDate   string                `json:"startDate"`
	EndDate     string                `json:"endDate"`
	Status      StatusResponse        `json:"status"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

func NewFullAdminReservationResponse(reservation *entity.Reservation) *FullAdminReservationResponse {
	return &FullAdminReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBriefBuildingResponse(reservation.Building),
		Tenant:      *NewBriefTenantResponse(reservation.User),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Status:      *NewStatusResponse(reservation.Status),
		CreatedAt:   reservation.CreatedAt,
		UpdatedAt:   reservation.UpdatedAt,
	}
}

type BriefTenantResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func NewBriefTenantResponse(user entity.User) *BriefTenantResponse {
	url := user.Detail.Picture.Url
	if url == "" {
		url = config.DEFAULT_USER_AVATAR
	}
	return &BriefTenantResponse{
		ID:      user.ID,
		Name:    user.Detail.Name,
		Email:   user.Email,
		Picture: url,
	}
}
