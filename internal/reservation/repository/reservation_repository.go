package repository

import (
	"office-booking-backend/pkg/entity"
	"time"

	"golang.org/x/net/context"
)

type ReservationRepository interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error)
	CountUserReservation(ctx context.Context, userID string) (int64, error)
	GetUserReservations(ctx context.Context, userID string, offset int, limit int) (*entity.Reservations, error)
	GetReservationByID(ctx context.Context, reservationID string) (*entity.Reservation, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error)
	AddBuildingReservation(ctx context.Context, reservation *entity.Reservation) error
	UpdateReservation(ctx context.Context, reservation *entity.Reservation) error
	DeleteReservationByID(ctx context.Context, reservationID string) error
}
