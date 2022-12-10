package impl

import (
	"context"
	"office-booking-backend/internal/review/repository"
	"office-booking-backend/pkg/entity"
	err2 "office-booking-backend/pkg/errors"

	"gorm.io/gorm"
)

type ReviewRepositoryImpl struct {
	db *gorm.DB
}

func NewReviewRepositoryImpl(db *gorm.DB) repository.ReviewRepository {
	return &ReviewRepositoryImpl{
		db: db,
	}
}

func (rv *ReviewRepositoryImpl) AddReservationReview(ctx context.Context, review *entity.Review) error {
	err := rv.db.WithContext(ctx).Create(review).Error
	if err != nil {
		return err
	}

	return nil
}

func (rv *ReviewRepositoryImpl) DeleteReviewByID(ctx context.Context, reservationID string) error {
	res := rv.db.WithContext(ctx).
		Model(&entity.Review{}).
		Where("id = ?", reservationID).
		Delete(&entity.Review{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return err2.ErrReservationNotFound
	}

	return nil
}
