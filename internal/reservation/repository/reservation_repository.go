package repository

import "golang.org/x/net/context"

type ReservationRepository interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error)
}
