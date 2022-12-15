package impl

import (
	"fmt"
	"log"
	repository2 "office-booking-backend/internal/building/repository"
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/internal/reservation/service"
	"office-booking-backend/pkg/custom"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"golang.org/x/net/context"
)

type ReservationServiceImpl struct {
	config       *viper.Viper
	repo         repository.ReservationRepository
	buildingRepo repository2.BuildingRepository
}

func NewReservationServiceImpl(reservationRepository repository.ReservationRepository, buildingRepository repository2.BuildingRepository, config *viper.Viper) service.ReservationService {
	return &ReservationServiceImpl{
		repo:         reservationRepository,
		buildingRepo: buildingRepository,
		config:       config,
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

	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		total, err := r.repo.GetReservationCountByTime(ctx)
		if err != nil {
			log.Println("error while getting this year reservation count: ", err)
			return err
		}

		statByTimeframe = dto.NewTimeframeStat(total)
		return nil
	})

	errGroup.Go(func() error {
		stat, err := r.repo.GetReservationCountByStatus(ctx)
		if err != nil {
			log.Println("error while getting total reservation count: ", err)
			return err
		}
		statByStatus = dto.NewReservationStatsResponse(stat)
		return nil
	})

	err := errGroup.Wait()
	if err != nil {
		return nil, err
	}

	res := &dto.ReservationStatResponse{
		ByTimeframe: statByTimeframe,
		ByStatus:    statByStatus,
	}

	return res, nil
}

func (r *ReservationServiceImpl) GetTotalRevenueByTime(ctx context.Context) (*dto.TimeframeStat, error) {
	total, err := r.repo.GetTotalRevenue(ctx)
	if err != nil {
		log.Println("error while getting this year revenue: ", err)
		return nil, err
	}

	res := dto.NewTimeframeStat(total)
	return res, nil
}

func (r *ReservationServiceImpl) CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error) {
	errGroup := errgroup.Group{}
	var building *entity.Building
	errGroup.Go(func() error {
		b, err := r.buildingRepo.GetBuildingDetailByID(ctx, reservation.BuildingID, true)
		if err != nil {
			if err == err2.ErrBuildingNotFound {
				return err2.ErrBuildingNotAvailable
			}

			log.Println("error while checking building availability: ", err)
			return err
		}
		building = b
		return nil
	})

	errGroup.Go(func() error {
		endDate := reservation.StartDate.ToTime().AddDate(0, reservation.Duration, 0)
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

	err := errGroup.Wait()
	if err != nil {
		return "", err
	}

	reservationEntity := reservation.ToEntity(userID)
	yearDuration := reservation.Duration / 12
	monthDuration := reservation.Duration - (yearDuration * 12)
	ammount := building.MonthlyPrice*monthDuration + building.AnnualPrice*yearDuration
	reservationEntity.Amount = ammount
	fmt.Printf("year: %d, month: %d, ammount: %d", yearDuration, monthDuration, ammount)
	err = r.repo.AddBuildingReservation(ctx, reservationEntity)
	if err != nil {
		log.Println("error while creating reservation: ", err)
		return "", err
	}
	return reservationEntity.ID, nil
}

func (r *ReservationServiceImpl) CreateAdminReservation(ctx context.Context, reservation *dto.AddAdminReservartionRequest) (string, error) {
	var building *entity.Building
	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		b, err := r.buildingRepo.GetBuildingDetailByID(ctx, reservation.BuildingID, true)
		if err != nil {
			if err == err2.ErrBuildingNotFound {
				return err2.ErrBuildingNotAvailable
			}

			log.Println("error while checking building availability: ", err)
			return err
		}
		building = b
		return nil
	})

	errGroup.Go(func() error {
		endDate := reservation.StartDate.ToTime().AddDate(0, reservation.Duration, 0)
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
	yearDuration := reservation.Duration / 12
	monthDuration := reservation.Duration - (yearDuration * 12)
	ammount := building.MonthlyPrice*monthDuration + building.AnnualPrice*yearDuration
	reservationEntity.Amount = ammount
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

	var building *entity.Building
	newReservation := reservation.ToEntity(reservationID)

	if reservation.BuildingID != "" || !reservation.StartDate.ToTime().IsZero() || reservation.Duration > 0 {
		buildingID := savedReservation.BuildingID
		if reservation.BuildingID != "" {
			buildingID = reservation.BuildingID
		}

		startDate := savedReservation.StartDate
		if !reservation.StartDate.ToTime().IsZero() {
			startDate = reservation.StartDate.ToTime()
		} else {
			// if start date is not provided, we will use the saved start date so we can calculate the end date correctly
			reservation.StartDate = custom.Date(savedReservation.StartDate)
		}

		endDate := savedReservation.EndDate
		if reservation.Duration != 0 {
			endDate = startDate.AddDate(0, reservation.Duration, 0)
		}

		errGroup := errgroup.Group{}
		errGroup.Go(func() error {
			b, err := r.buildingRepo.GetBuildingDetailByID(ctx, reservation.BuildingID, true)
			if err != nil {
				if err == err2.ErrBuildingNotFound {
					return err2.ErrBuildingNotAvailable
				}

				log.Println("error while checking building availability: ", err)
				return err
			}
			building = b
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

		yearDuration := reservation.Duration / 12
		monthDuration := reservation.Duration - (yearDuration * 12)
		ammount := building.MonthlyPrice*monthDuration + building.AnnualPrice*yearDuration
		newReservation.Amount = ammount
	}

	err = r.repo.UpdateReservation(ctx, newReservation)
	if err != nil {
		log.Println("error while updating reservation: ", err)
		return err
	}

	return nil
}

func (r *ReservationServiceImpl) UpdateReservationStatus(ctx context.Context, reservationID string, statusRequest *dto.UpdateReservationStatusRequest) error {
	// 1 = pending, 2 = rejected, 3 = cancelled, 4 = awaiting payment, 5 = active, 6 = completed
	reservationEntity := statusRequest.ToEntity(reservationID)
	if statusRequest.StatusID == 4 {
		reservationEntity.AcceptedAt = time.Now()
		reservationEntity.ExpiredAt = time.Now().Add(r.config.GetDuration("reservation.expiredAt"))
	}

	err := r.repo.UpdateReservation(ctx, reservationEntity)
	if err != nil {
		log.Println("error while updating reservation status: ", err)
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
