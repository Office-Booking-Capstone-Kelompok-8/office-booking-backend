package repository

import (
	"golang.org/x/net/context"
	"time"
)

type ReservationRepository interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error)
}
