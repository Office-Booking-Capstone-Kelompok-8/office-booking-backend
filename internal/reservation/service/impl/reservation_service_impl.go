package impl

import (
	"golang.org/x/net/context"
	"log"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/internal/reservation/service"
)

type ReservationServiceImpl struct {
	repo repository.ReservationRepository
}

func NewReservationServiceImpl(repo repository.ReservationRepository) service.ReservationService {
	return &ReservationServiceImpl{
		repo: repo,
	}
}

func (r *ReservationServiceImpl) CountUserActiveReservations(ctx context.Context, userID string) (int64, error) {
	count, err := r.repo.CountUserActiveReservations(ctx, userID)
	if err != nil {
		log.Println("error while counting user reservations: ", err)
		return 0, err
	}

	return count, nil
}
