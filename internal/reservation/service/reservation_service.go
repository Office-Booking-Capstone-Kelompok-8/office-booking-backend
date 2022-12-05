package service

import (
	"office-booking-backend/internal/reservation/dto"
	"time"

	"golang.org/x/net/context"
)

type ReservationService interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error)
	GetUserReservations(ctx context.Context, userID string, page int, limit int) (*dto.BriefReservationsResponse, int64, error)
	GetReservationByID(ctx context.Context, reservationID string) (*dto.FullAdminReservationResponse, error)
	CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error)
	CreateAdminReservation(ctx context.Context, reservation *dto.AddAdminReservartionRequest) (string, error)
	CancelReservation(ctx context.Context, userID string, reservationID string) error
	DeleteReservationByID(ctx context.Context, reservationID string) error
}
