package repository

import (
	"office-booking-backend/pkg/entity"
	"time"

	"golang.org/x/net/context"
)

type ReservationRepository interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error)
	GetUserReservations(ctx context.Context, userID string, offset int, limit int) (*entity.Reservations, int64, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error)
	AddBuildingReservation(ctx context.Context, reservation *entity.Reservation) error
}
