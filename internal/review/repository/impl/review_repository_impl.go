package impl

import (
	"context"
	"office-booking-backend/internal/review/repository"
	"office-booking-backend/pkg/entity"

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
