package impl

import (
	"log"
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/internal/reservation/service"
	err2 "office-booking-backend/pkg/errors"
	"time"

	"golang.org/x/net/context"
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

func NewReservationsServiceImpl(repo repository.ReservationRepository) service.ReservationsService {
	return &ReservationServiceImpl{
		repo: repo,
	}
}

func (r *ReservationServiceImpl) GetAllReservations(ctx context.Context, status string, buildingID string, userID string, userName string, startDate time.Time, endDate time.Time, page int, limit int) (*dto.BriefReservationResponse, int64, error) {

	if startDate.After(endDate) {
		return nil, 0, err2.ErrStartDateAfterEndDate
	}

	offset := (page - 1) * limit

	reservations, count, err := r.repo.GetAllReservations(ctx, status, buildingID, userID, userName, startDate, endDate, offset, limit)
	if err != nil {
		log.Println("error when getting all reservations: ", err)
		return nil, 0, err
	}

	reservationsResponse := dto.NewBriefReservationResponse(reservations)
	return reservationsResponse, count, nil
}
