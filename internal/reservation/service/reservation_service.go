package service

import (
	"office-booking-backend/internal/reservation/dto"
	"time"

	"golang.org/x/net/context"
)

type ReservationService interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	GetUserReservations(ctx context.Context, userID string, page int, limit int) (*dto.BriefReservationsResponse, int64, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error)
	CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error)
}
