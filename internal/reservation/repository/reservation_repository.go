package repository

import (
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/pkg/entity"
	"time"

	"golang.org/x/net/context"
)

type ReservationRepository interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error)
	CountUserReservation(ctx context.Context, userID string) (int64, error)
	CountReservation(ctx context.Context, filter *dto.ReservationQueryParam) (int64, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time, excludedReservationID ...string) (bool, error)
	GetReservations(ctx context.Context, filter *dto.ReservationQueryParam) (*entity.Reservations, error)
	GetUserReservations(ctx context.Context, userID string, offset int, limit int) (*entity.Reservations, error)
	GetReservationByID(ctx context.Context, reservationID string) (*entity.Reservation, error)
	GetUserReservationByID(ctx context.Context, reservationID string, userID string) (*entity.Reservation, error)
	GetReservationReview(ctx context.Context, reservations *entity.Reservation) (*entity.Review, error)
	GetReservationCountByStatus(ctx context.Context) (*entity.StatusesStat, error)
	GetReservationCountByTime(ctx context.Context) (*entity.TimeframeStat, error)
	GetTotalRevenue(ctx context.Context) (*entity.TimeframeStat, error)
	GetReservationTaskUntilToday(ctx context.Context) (*entity.Reservations, error)
	AddBuildingReservation(ctx context.Context, reservation *entity.Reservation) error
	AddReservationReviews(ctx context.Context, review *entity.Review) error
	UpdateReservation(ctx context.Context, reservation *entity.Reservation) error
	UpdateReservationReviews(ctx context.Context, review *entity.Review) error
	DeleteReservationByID(ctx context.Context, reservationID string) error
}
