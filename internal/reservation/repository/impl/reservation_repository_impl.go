package impl

import (
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"
	"time"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ReservationRepositoryImpl struct {
	db *gorm.DB
}

func NewReservationRepositoryImpl(db *gorm.DB) repository.ReservationRepository {
	return &ReservationRepositoryImpl{
		db: db,
	}
}

func (r *ReservationRepositoryImpl) CountUserActiveReservations(ctx context.Context, userID string) (int64, error) {
	var count int64
	// Count user active reservations with status id not 2 (rejected), 3 (canceled) or 5 (completed) and not expired
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("user_id = ? AND status_id NOT IN (2, 3, 5) AND end_date > ?", userID, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error) {
	var count int64
	// Count building active reservations with status id not 2 (rejected), 3 (canceled) or 5 (completed) and not expired
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("building_id = ? AND status_id NOT IN (2, 3, 5) AND end_date > ?", buildingID, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error) {
	var count int64
	// Count building active reservations with status id not 2 (rejected), 3 (canceled) or 5 (completed)
	// and not in the same time range as the new reservation
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("building_id = ? AND status_id NOT IN (2, 3, 5)", buildingID).
		Where("start_date <= ? AND end_date >= ?", end, start).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *ReservationRepositoryImpl) GetUserReservations(ctx context.Context, userID string, offset int, limit int) (*entity.Reservations, int64, error) {
	var count int64
	var reservations entity.Reservations

	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Joins("Status").
		Joins("Building").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&reservations).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return &reservations, count, nil
}

func (r *ReservationRepositoryImpl) GetReservationByID(ctx context.Context, reservationID string) (*entity.Reservation, error) {
	var reservation entity.Reservation

	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Joins("Status").
		Joins("Building").
		Where("`reservations`.`id` = ?", reservationID).
		First(&reservation).
		Error
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *ReservationRepositoryImpl) AddBuildingReservation(ctx context.Context, reservation *entity.Reservation) error {
	err := r.db.WithContext(ctx).Create(reservation).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ReservationRepositoryImpl) UpdateReservation(ctx context.Context, reservation *entity.Reservation) error {
	res := r.db.WithContext(ctx).
		Model(entity.Reservation{}).
		Where("id = ?", reservation.ID).
		Updates(reservation)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrReservationNotFound
	}

	return nil
}

func (r *ReservationRepositoryImpl) DeleteReservationByID(ctx context.Context, reservationID string) error {
	res := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("id = ?", reservationID).
		Delete(&entity.Reservation{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrReservationNotFound
	}

	return nil
}
