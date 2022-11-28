package impl

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/entity"
	"time"
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
