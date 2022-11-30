package repository

import (
	"office-booking-backend/pkg/entity"
	"time"

	"golang.org/x/net/context"
)

type ReservationRepository interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	GetAllReservations(ctx context.Context, status string, buildingID string, userID string, userName string, startDate time.Time, endDate time.Time, limit int, offset int) (*entity.Reservation, int64, error)
}
