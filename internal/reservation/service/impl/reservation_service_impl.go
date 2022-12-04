package impl

import (
	"golang.org/x/sync/errgroup"
	"log"
	repository2 "office-booking-backend/internal/building/repository"
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/internal/reservation/service"
	err2 "office-booking-backend/pkg/errors"
	"time"

	"golang.org/x/net/context"
)

type ReservationServiceImpl struct {
	repo         repository.ReservationRepository
	buildingRepo repository2.BuildingRepository
}

func NewReservationServiceImpl(reservationRepository repository.ReservationRepository, buildingRepository repository2.BuildingRepository) service.ReservationService {
	return &ReservationServiceImpl{
		repo:         reservationRepository,
		buildingRepo: buildingRepository,
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

func (r *ReservationServiceImpl) IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error) {
	isAvailable, err := r.repo.IsBuildingAvailable(ctx, buildingID, start, end)
	if err != nil {
		log.Println("error while checking building availability: ", err)
		return false, err
	}

	return isAvailable, nil
}

func (r *ReservationServiceImpl) GetUserReservations(ctx context.Context, userID string, page int, limit int) (*dto.BriefReservationsResponse, int64, error) {
	offset := (page - 1) * limit
	reservations, count, err := r.repo.GetUserReservations(ctx, userID, offset, limit)
	if err != nil {
		log.Println("error while getting user reservations: ", err)
		return nil, 0, err
	}

	reservationDto := dto.NewBriefReservationsResponse(reservations)
	return reservationDto, count, nil
}

func (r *ReservationServiceImpl) CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error) {
	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		isPublished, err := r.buildingRepo.IsBuildingPublished(ctx, reservation.BuildingID)
		if err != nil {
			log.Println("error while checking building availability: ", err)
			return err
		}

		if !isPublished {
			return err2.ErrBuildingNotAvailable
		}

		return nil
	})

	errGroup.Go(func() error {
		isAvailable, err := r.repo.IsBuildingAvailable(ctx, reservation.BuildingID, reservation.StartDate.ToTime(), reservation.EndDate.ToTime())
		if err != nil {
			log.Println("error while checking building availability: ", err)
			return err
		}

		if !isAvailable {
			return err2.ErrBuildingNotAvailable
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		return "", err
	}

	reservationEntity := reservation.ToEntity(userID)
	err := r.repo.AddBuildingReservation(ctx, reservationEntity)
	if err != nil {
		log.Println("error while creating reservation: ", err)
		return "", err
	}

	return reservationEntity.ID, nil
}
