package dto

import (
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/entity"
	"time"
)

func newDuration(start, end time.Time) int {
	yearStart, monthStart, _ := start.Date()
	yearEnd, monthEnd, _ := end.Date()
	return (yearEnd-yearStart)*12 + (int(monthEnd) - int(monthStart))
}

type BriefReservationResponse struct {
	ID          string                `json:"id"`
	Building    BriefBuildingResponse `json:"building"`
	CompanyName string                `json:"companyName"`
	StartDate   string                `json:"startDate"`
	EndDate     string                `json:"endDate"`
	Duration    int                   `json:"duration"`
	Amount      int                   `json:"amount"`
	Status      StatusResponse        `json:"status"`
	ExpiredAt   string                `json:"expiredAt"`
	CreatedAt   string                `json:"createdAt"`
}

func NewBriefReservationResponse(reservation *entity.Reservation) *BriefReservationResponse {
	return &BriefReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBriefBuildingResponse(reservation.Building),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Duration:    newDuration(reservation.StartDate, reservation.EndDate),
		Amount:      reservation.Amount,
		Status:      *NewStatusResponse(reservation.Status),
		ExpiredAt:   reservation.ExpiredAt.Format(config.DATE_RESPONSE_FORMAT),
		CreatedAt:   reservation.CreatedAt.Format(config.DATE_RESPONSE_FORMAT),
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

type BriefAdminReservationResponse struct {
	ID          string                `json:"id"`
	Building    BriefBuildingResponse `json:"building"`
	Tenant      TenantResponse        `json:"tenant"`
	CompanyName string                `json:"companyName"`
	StartDate   string                `json:"startDate"`
	EndDate     string                `json:"endDate"`
	Duration    int                   `json:"duration"`
	Amount      int                   `json:"amount"`
	Status      StatusResponse        `json:"status"`
	ExpiredAt   string                `json:"expiredAt"`
	CreatedAt   string                `json:"createdAt"`
}

func NewBriefAdminReservationResponse(reservation *entity.Reservation) *BriefAdminReservationResponse {
	return &BriefAdminReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBriefBuildingResponse(reservation.Building),
		Tenant:      *NewTenantResponse(reservation.User),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Duration:    newDuration(reservation.StartDate, reservation.EndDate),
		Amount:      reservation.Amount,
		Status:      *NewStatusResponse(reservation.Status),
		ExpiredAt:   reservation.ExpiredAt.Format(config.DATE_RESPONSE_FORMAT),
		CreatedAt:   reservation.CreatedAt.Format(config.DATE_RESPONSE_FORMAT),
	}
}

type BriefAdminReservationsResponse []BriefAdminReservationResponse

func NewBriefAdminReservationsResponse(reservations *entity.Reservations) *BriefAdminReservationsResponse {
	response := new(BriefAdminReservationsResponse)
	for _, reservation := range *reservations {
		*response = append(*response, *NewBriefAdminReservationResponse(&reservation))
	}
	return response
}

type BuildingResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func NewBuildingResponse(building entity.Building) *BuildingResponse {
	return &BuildingResponse{
		ID:       building.ID,
		Name:     building.Name,
		Picture:  building.Pictures[0].ThumbnailUrl,
		City:     building.City.Name,
		District: building.District.Name,
		Address:  building.Address,
	}
}

type BriefBuildingResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	City    string `json:"city"`
}

func NewBriefBuildingResponse(building entity.Building) *BriefBuildingResponse {
	return &BriefBuildingResponse{
		ID:      building.ID,
		Name:    building.Name,
		Picture: building.Pictures[0].ThumbnailUrl,
		City:    building.City.Name,
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
	ID          string           `json:"id"`
	Building    BuildingResponse `json:"building"`
	Tenant      TenantResponse   `json:"tenant"`
	CompanyName string           `json:"companyName"`
	StartDate   string           `json:"startDate"`
	EndDate     string           `json:"endDate"`
	Duration    int              `json:"duration"`
	Amount      int              `json:"amount"`
	Status      StatusResponse   `json:"status"`
	Message     string           `json:"message"`
	ExpiredAt   string           `json:"expiredAt"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

func NewFullAdminReservationResponse(reservation *entity.Reservation) *FullAdminReservationResponse {
	return &FullAdminReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBuildingResponse(reservation.Building),
		Tenant:      *NewTenantResponse(reservation.User),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Duration:    newDuration(reservation.StartDate, reservation.EndDate),
		Amount:      reservation.Amount,
		Status:      *NewStatusResponse(reservation.Status),
		Message:     reservation.Message,
		ExpiredAt:   reservation.ExpiredAt.Format(config.DATE_RESPONSE_FORMAT),
		CreatedAt:   reservation.CreatedAt.Format(config.DATE_RESPONSE_FORMAT),
		UpdatedAt:   reservation.UpdatedAt.Format(config.DATE_RESPONSE_FORMAT),
	}
}

type TenantResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func NewTenantResponse(user entity.User) *TenantResponse {
	url := user.Detail.Picture.Url
	if url == "" {
		url = config.DEFAULT_USER_AVATAR
	}
	return &TenantResponse{
		ID:      user.ID,
		Name:    user.Detail.Name,
		Email:   user.Email,
		Picture: url,
	}
}

type FullReservationResponse struct {
	ID          string           `json:"id"`
	Building    BuildingResponse `json:"building"`
	CompanyName string           `json:"companyName"`
	StartDate   string           `json:"startDate"`
	EndDate     string           `json:"endDate"`
	Duration    int              `json:"duration"`
	Status      StatusResponse   `json:"status"`
	Amount      int              `json:"amount"`
	Message     string           `json:"message"`
	ExpiredAt   string           `json:"expiredAt"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

func NewFullReservationResponse(reservation *entity.Reservation) *FullReservationResponse {
	return &FullReservationResponse{
		ID:          reservation.ID,
		Building:    *NewBuildingResponse(reservation.Building),
		CompanyName: reservation.CompanyName,
		StartDate:   reservation.StartDate.Format(config.DATE_RESPONSE_FORMAT),
		EndDate:     reservation.EndDate.Format(config.DATE_RESPONSE_FORMAT),
		Duration:    newDuration(reservation.StartDate, reservation.EndDate),
		Amount:      reservation.Amount,
		Status:      *NewStatusResponse(reservation.Status),
		Message:     reservation.Message,
		ExpiredAt:   reservation.ExpiredAt.Format(config.DATE_RESPONSE_FORMAT),
		CreatedAt:   reservation.CreatedAt.Format(config.DATE_RESPONSE_FORMAT),
		UpdatedAt:   reservation.UpdatedAt.Format(config.DATE_RESPONSE_FORMAT),
	}
}

type ReservationStatResponse struct {
	ByStatus    *ReservationTotals `json:"byStatus"`
	ByTimeframe *TimeframeStat     `json:"byTimeframe"`
}

type ReservationTotal struct {
	StatusID   int    `json:"statusId"`
	StatusName string `json:"statusName"`
	Total      int    `json:"count"`
}

func NewReservationTotal(stat entity.StatusStat) *ReservationTotal {
	return &ReservationTotal{
		StatusID:   int(stat.StatusID),
		StatusName: stat.StatusName,
		Total:      int(stat.Total),
	}
}

type ReservationTotals []ReservationTotal

func NewReservationStatsResponse(stats *entity.StatusesStat) *ReservationTotals {
	response := new(ReservationTotals)
	for _, stat := range *stats {
		*response = append(*response, *NewReservationTotal(stat))
	}
	return response
}

type TimeframeStat struct {
	Today     int64 `json:"today"`
	ThisWeek  int64 `json:"thisWeek"`
	ThisMonth int64 `json:"thisMonth"`
	ThisYear  int64 `json:"thisYear"`
	AllTime   int64 `json:"allTime"`
}

func NewTimeframeStat(stats *entity.TimeframeStat) *TimeframeStat {
	return &TimeframeStat{
		Today:     stats.Day,
		ThisWeek:  stats.Week,
		ThisMonth: stats.Month,
		ThisYear:  stats.Year,
		AllTime:   stats.All,
	}
}
