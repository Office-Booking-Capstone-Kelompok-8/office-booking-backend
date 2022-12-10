package service

import (
	"context"
	"office-booking-backend/internal/review/dto"
)

type ReviewService interface {
	CreateReservationReview(ctx context.Context, review *dto.AddReviewRequest) error
}
