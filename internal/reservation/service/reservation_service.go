package service

import (
	"office-booking-backend/internal/reservation/dto"
	"time"

	"golang.org/x/net/context"
)

type ReservationService interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
}

type ReservationsService interface {
	GetAllReservations(ctx context.Context, status string, buildingID string, userID string, userName string, startDate time.Time, endDate time.Time, page int, limit int) (*dto.BriefReservationResponse, int64, error)
}
