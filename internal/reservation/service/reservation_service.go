package service

import (
	"office-booking-backend/internal/reservation/dto"
	"time"

	"golang.org/x/net/context"
)

type ReservationService interface {
	CountUserActiveReservations(ctx context.Context, userID string) (int64, error)
	GetReservationByID(ctx context.Context, reservationID string) (*dto.FullAdminReservationResponse, error)
	GetReservations(ctx context.Context, filter *dto.ReservationQueryParam) (*dto.BriefAdminReservationsResponse, int64, error)
	GetUserReservationByID(ctx context.Context, userID string, reservationID string) (*dto.FullReservationResponse, error)
	GetUserReservations(ctx context.Context, userID string, page int, limit int) (*dto.BriefReservationsResponse, int64, error)
	GetReservationStat(ctx context.Context) (*dto.ReservationStatResponse, error)
	GetTotalRevenueByTime(ctx context.Context) (*dto.TimeframeStat, error)
	GetReservationReview(ctx context.Context, reservationID string, userID string) (*dto.BriefReviewResponse, error)
	IsBuildingAvailable(ctx context.Context, buildingID string, startDate time.Time, duration int) (bool, error)
	CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error)
	CreateAdminReservation(ctx context.Context, reservation *dto.AddAdminReservartionRequest) (string, error)
	CreateReservationReviews(ctx context.Context, review *dto.AddReviewRequest) error
	CancelReservation(ctx context.Context, userID string, reservationID string) error
	UpdateReservation(ctx context.Context, reservationID string, reservation *dto.UpdateReservationRequest) error
	UpdateReservationStatus(ctx context.Context, reservationID string, statusRequest *dto.UpdateReservationStatusRequest) error
	DeleteReservationByID(ctx context.Context, reservationID string) error
}
