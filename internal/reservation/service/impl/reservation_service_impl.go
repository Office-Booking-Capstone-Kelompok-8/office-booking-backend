package impl

import (
	"log"
	repository2 "office-booking-backend/internal/building/repository"
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/internal/reservation/service"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"time"

	"golang.org/x/sync/errgroup"

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

func (r *ReservationServiceImpl) IsBuildingAvailable(ctx context.Context, buildingID string, startDate time.Time, duration int) (bool, error) {
	endDate := startDate.AddDate(0, duration, 0)
	isAvailable, err := r.repo.IsBuildingAvailable(ctx, buildingID, startDate, endDate)
	if err != nil {
		log.Println("error while checking building availability: ", err)
		return false, err
	}

	return isAvailable, nil
}

func (r *ReservationServiceImpl) GetUserReservations(ctx context.Context, userID string, page int, limit int) (*dto.BriefReservationsResponse, int64, error) {
	count, err := r.repo.CountUserReservation(ctx, userID)
	if err != nil {
		log.Println("error while counting user reservations: ", err)
		return nil, 0, err
	}

	if count == 0 {
		return nil, 0, nil
	}

	offset := (page - 1) * limit
	reservations, err := r.repo.GetUserReservations(ctx, userID, offset, limit)
	if err != nil {
		log.Println("error while getting user reservations: ", err)
		return nil, 0, err
	}

	reservationDto := dto.NewBriefReservationsResponse(reservations)
	return reservationDto, count, nil
}

func (r *ReservationServiceImpl) GetReservations(ctx context.Context, filter *dto.ReservationQueryParam) (*dto.BriefAdminReservationsResponse, int64, error) {
	count, err := r.repo.CountReservation(ctx, filter)
	if err != nil {
		log.Println("error while counting reservations: ", err)
		return nil, 0, err
	}

	if count == 0 {
		return nil, 0, nil
	}

	filter.Offset = (filter.Page - 1) * filter.Limit
	reservations, err := r.repo.GetReservations(ctx, filter)
	if err != nil {
		log.Println("error while getting reservations: ", err)
		return nil, 0, err
	}

	reservationDto := dto.NewBriefAdminReservationsResponse(reservations)
	return reservationDto, count, nil
}

func (r *ReservationServiceImpl) GetUserReservationByID(ctx context.Context, reservationID string, userID string) (*dto.FullReservationResponse, error) {
	reservation, err := r.repo.GetUserReservationByID(ctx, reservationID, userID)
	if err != nil {
		log.Println("error while getting reservation by id: ", err)
		return nil, err
	}

	reservationDto := dto.NewFullReservationResponse(reservation)
	return reservationDto, nil
}

func (r *ReservationServiceImpl) GetReservationByID(ctx context.Context, reservationID string) (*dto.FullAdminReservationResponse, error) {
	reservation, err := r.repo.GetReservationByID(ctx, reservationID)
	if err != nil {
		log.Println("error while getting reservation by id: ", err)
		return nil, err
	}

	reservationDto := dto.NewFullAdminReservationResponse(reservation)
	return reservationDto, nil
}

func (r *ReservationServiceImpl) GetReservationStat(ctx context.Context) (*dto.ReservationStatResponse, error) {
	statByStatus := new(dto.ReservationTotals)
	statByTimeframe := new(dto.TimeframeStat)

	errgroup := errgroup.Group{}

	errgroup.Go(func() error {
		total, err := r.repo.GetReservationCount(ctx)
		if err != nil {
			log.Println("error while getting this year reservation count: ", err)
			return err
		}

		statByTimeframe = dto.NewTimeframeStat(total)
		return nil
	})

	errgroup.Go(func() error {
		stat, err := r.repo.GetReservationTotal(ctx)
		if err != nil {
			log.Println("error while getting total reservation count: ", err)
			return err
		}
		statByStatus = dto.NewReservationStatsResponse(stat)
		return nil
	})

	err := errgroup.Wait()
	if err != nil {
		return nil, err
	}

	res := &dto.ReservationStatResponse{
		ByTimeframe: statByTimeframe,
		ByStatus:    statByStatus,
	}

	return res, nil
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
		endDate := reservation.StartDate.AddDate(0, reservation.Duration, 0)
		isAvailable, err := r.repo.IsBuildingAvailable(ctx, reservation.BuildingID, reservation.StartDate.ToTime(), endDate)
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

func (r *ReservationServiceImpl) CreateAdminReservation(ctx context.Context, reservation *dto.AddAdminReservartionRequest) (string, error) {
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
		endDate := reservation.StartDate.AddDate(0, reservation.Duration, 0)
		isAvailable, err := r.repo.IsBuildingAvailable(ctx, reservation.BuildingID, reservation.StartDate.ToTime(), endDate)
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

	reservationEntity := reservation.ToEntity()
	err := r.repo.AddBuildingReservation(ctx, reservationEntity)
	if err != nil {
		log.Println("error while creating reservation: ", err)
		return "", err
	}

	return reservationEntity.ID, nil
}

func (r *ReservationServiceImpl) CancelReservation(ctx context.Context, userID string, reservationID string) error {
	reservation, err := r.repo.GetReservationByID(ctx, reservationID)
	if err != nil {
		log.Println("error while getting reservation by id: ", err)
		return err
	}

	if reservation == nil {
		return err2.ErrReservationNotFound
	}

	if reservation.UserID != userID {
		return err2.ErrNoPermission
	}

	if reservation.StatusID == 5 {
		return err2.ErrReservationActive
	}

	newReservation := &entity.Reservation{
		ID:       reservationID,
		StatusID: 3,
	}

	err = r.repo.UpdateReservation(ctx, newReservation)
	if err != nil {
		log.Println("error while updating reservation: ", err)
		return err
	}

	return nil
}

func (r *ReservationServiceImpl) UpdateReservation(ctx context.Context, reservationID string, reservation *dto.UpdateReservationRequest) error {
	savedReservation, err := r.repo.GetReservationByID(ctx, reservationID)
	if err != nil {
		log.Println("error while getting reservation by id: ", err)
		return err
	}
	if savedReservation == nil {
		return err2.ErrReservationNotFound
	}

	if reservation.BuildingID != "" || !reservation.StartDate.IsZero() || reservation.Duration > 0 {
		buildingID := savedReservation.BuildingID
		if reservation.BuildingID != "" {
			buildingID = reservation.BuildingID
		}

		startDate := savedReservation.StartDate
		if !reservation.StartDate.IsZero() {
			startDate = reservation.StartDate.ToTime()
		} else {
			// if start date is not provided, we will use the saved start date so we can calculate the end date correctly
			reservation.StartDate.Time = savedReservation.StartDate
		}

		endDate := savedReservation.EndDate
		if reservation.Duration != 0 {
			endDate = startDate.AddDate(0, reservation.Duration, 0)
		}

		errGroup := errgroup.Group{}
		errGroup.Go(func() error {
			isPublished, err := r.buildingRepo.IsBuildingPublished(ctx, buildingID)
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
			isAvailable, err := r.repo.IsBuildingAvailable(ctx, buildingID, startDate, endDate, reservationID)
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
			return err
		}
	}

	newReservation := reservation.ToEntity(reservationID)
	err = r.repo.UpdateReservation(ctx, newReservation)
	if err != nil {
		log.Println("error while updating reservation: ", err)
		return err
	}

	return nil
}

func (r *ReservationServiceImpl) DeleteReservationByID(ctx context.Context, reservationID string) error {
	err := r.repo.DeleteReservationByID(ctx, reservationID)
	if err != nil {
		log.Println("error while deleting reservation by id: ", err)
		return err
	}

	return nil
}

// get review for reservation
func (r *ReservationServiceImpl) GetReservationReviews(ctx context.Context) (*dto.BriefReviewsResponse, error) {
	reviews, err := r.repo.GetReservationReviews(ctx)
	if err != nil {
		log.Println("error when getting all banks: ", err)
		return nil, err
	}

	return dto.NewBriefReviewsResponse(reviews), nil
}
