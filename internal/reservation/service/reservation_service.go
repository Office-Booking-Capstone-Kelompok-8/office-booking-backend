package service

import "golang.org/x/net/context"

type ReservationService interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
}
