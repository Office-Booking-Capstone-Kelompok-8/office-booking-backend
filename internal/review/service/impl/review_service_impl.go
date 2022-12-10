package impl

import (
	"context"
	"log"
	"office-booking-backend/internal/review/dto"
	"office-booking-backend/internal/review/repository"
	"office-booking-backend/internal/review/service"
)

type ReviewServiceImpl struct {
	repo repository.ReviewRepository
}

func NewReviewServiceImpl(reviewRepository repository.ReviewRepository) service.ReviewService {
	return &ReviewServiceImpl{
		repo: reviewRepository,
	}
}

func (rv *ReviewServiceImpl) CreateReservationReview(ctx context.Context, review *dto.AddReviewRequest) error {
	reviewEntity := review.ToEntity()
	err := rv.repo.AddReservationReview(ctx, reviewEntity)
	if err != nil {
		log.Println("error when create review: ", err)
		return err
	}

	return nil
}

func (rv *ReviewServiceImpl) DeleteReviewByID(ctx context.Context, reservationID string) error {
	err := rv.repo.DeleteReviewByID(ctx, reservationID)
	if err != nil {
		log.Println("error while deleting review by reservation id: ", err)
		return err
	}

	return nil
}
