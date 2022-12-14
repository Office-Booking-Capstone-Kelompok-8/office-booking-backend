package mock

import (
	"office-booking-backend/internal/reservation/dto"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type ReservationServiceMock struct {
	mock.Mock
}

func (r *ReservationServiceMock) GetReservationByID(ctx context.Context, reservationID string) (*dto.FullAdminReservationResponse, error) {
	args := r.Called(ctx, reservationID)
	return args.Get(0).(*dto.FullAdminReservationResponse), args.Error(1)
}

func (r *ReservationServiceMock) CountUserActiveReservations(ctx context.Context, userID string) (int64, error) {
	args := r.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ReservationServiceMock) GetUserReservations(ctx context.Context, userID string, page int, limit int) (*dto.BriefReservationsResponse, int64, error) {
	args := r.Called(ctx, userID, page, limit)
	return args.Get(0).(*dto.BriefReservationsResponse), args.Get(1).(int64), args.Error(2)
}

func (r *ReservationServiceMock) IsBuildingAvailable(ctx context.Context, buildingID string, startDate time.Time, duration int) (bool, error) {
	args := r.Called(ctx, buildingID, startDate, duration)
	return args.Get(0).(bool), args.Error(1)
}

func (r *ReservationServiceMock) CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error) {
	args := r.Called(ctx, userID, reservation)
	return args.Get(0).(string), args.Error(1)
}

func (r *ReservationServiceMock) CreateAdminReservation(ctx context.Context, reservation *dto.AddAdminReservartionRequest) (string, error) {
	arga := r.Called(ctx, reservation)
	return arga.Get(0).(string), arga.Error(1)
}

func (r *ReservationServiceMock) DeleteReservationByID(ctx context.Context, reservationID string) error {
	args := r.Called(ctx, reservationID)
	return args.Error(0)
}

func (r *ReservationServiceMock) CancelReservation(ctx context.Context, userID string, reservationID string) error {
	args := r.Called(ctx, userID, reservationID)
	return args.Error(0)
}

func (r *ReservationServiceMock) UpdateReservation(ctx context.Context, reservationID string, reservation *dto.UpdateReservationRequest) error {
	args := r.Called(ctx, reservationID, reservation)
	return args.Error(0)
}

func (r *ReservationServiceMock) GetUserReservationByID(ctx context.Context, userID string, reservationID string) (*dto.FullReservationResponse, error) {
	args := r.Called(ctx, userID, reservationID)
	return args.Get(0).(*dto.FullReservationResponse), args.Error(1)
}

func (r *ReservationServiceMock) GetReservations(ctx context.Context, filter *dto.ReservationQueryParam) (*dto.BriefAdminReservationsResponse, int64, error) {
	args := r.Called(ctx, filter)
	return args.Get(0).(*dto.BriefAdminReservationsResponse), args.Get(1).(int64), args.Error(2)
}

func (r *ReservationServiceMock) GetReservationStat(ctx context.Context) (*dto.ReservationStatResponse, error) {
	args := r.Called(ctx)
	return args.Get(0).(*dto.ReservationStatResponse), args.Error(1)
}

func (r *ReservationServiceMock) GetTotalRevenueByTime(ctx context.Context) (*dto.TimeframeStat, error) {
	args := r.Called(ctx)
	return args.Get(0).(*dto.TimeframeStat), args.Error(1)
}

func (r *ReservationServiceMock) UpdateReservationStatus(ctx context.Context, reservationID string, statusRequest *dto.UpdateReservationStatusRequest) error {
	args := r.Called(ctx, reservationID, statusRequest)
	return args.Error(0)
}

func (r *ReservationServiceMock) GetReservationReview(ctx context.Context, reservationID string, userID string) (*dto.BriefReviewResponse, error) {
	args := r.Called(ctx, reservationID, userID)
	return args.Get(0).(*dto.BriefReviewResponse), args.Error(1)
}

func (r *ReservationServiceMock) CreateReservationReview(ctx context.Context, review *dto.AddReviewRequest, reservationID string, userID string) error {
	args := r.Called(ctx, review, reservationID, userID)
	return args.Error(0)
}

func (r *ReservationServiceMock) UpdateReservationReview(ctx context.Context, review *dto.UpdateReviewRequest, reservationID string, userID string) error {
	args := r.Called(ctx, review, reservationID, userID)
	return args.Error(0)
}
