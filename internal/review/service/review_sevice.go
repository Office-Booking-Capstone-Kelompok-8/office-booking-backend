package service

import (
	"context"
	"office-booking-backend/internal/review/dto"
)

type ReviewService interface {
	CreateReservationReview(ctx context.Context, review *dto.AddReviewRequest) error
	DeleteReviewByID(ctx context.Context, reservationID string) error
}
