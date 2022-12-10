package repository

import (
	"context"
	"office-booking-backend/pkg/entity"
)

type ReviewRepository interface {
	AddReservationReview(ctx context.Context, review *entity.Review) error
	DeleteReviewByID(ctx context.Context, reservationID string) error
}
