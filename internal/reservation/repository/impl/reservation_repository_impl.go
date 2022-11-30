package impl

import (
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/entity"
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
	// Count user active reservations with status id not 3 (rejected), 4 (canceled) or 6 (completed) and not expired
	err := r.db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("user_id = ? AND status_id NOT IN (3, 4, 6) AND end_date > ?", userID, time.Now()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReservationRepositoryImpl) GetAllReservations(ctx context.Context, status string, buildingID string,
	userID string, userName string, startDate time.Time, endDate time.Time, limit int, offset int) (*entity.Reservation, int64, error) {
	reservations := &entity.Reservation{}
	var count int64

	query := r.db.WithContext(ctx).
		Joins("Building").
		Joins("User").Model(&entity.Reservation{})
	if status != "" {
		query = query.Where("`reservations`.`status` LIKE ?", "%"+status+"%")
	}

	if buildingID != "" {
		query = query.Where("`Building`.`id` = ?", buildingID)
	}

	if userID != "" {
		query = query.Where("`User`.`id` = ?", userID)
	}

	if userName != "" {
		query = query.Where("`User`.`name` = ?", userName)
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("NOT EXISTS (SELECT * FROM `reservations` WHERE `reservations`.`building_id` = `buildings`.`id` AND `reservations`.`start_date` <= ? AND `reservations`.`end_date` >= ?)", endDate, startDate)
	}

	err := query.Limit(limit).Offset(offset).Find(reservations).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return reservations, count, nil
}
