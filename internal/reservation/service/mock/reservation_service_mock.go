package mock

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"office-booking-backend/internal/reservation/dto"
	"time"
)

type ReservationServiceMock struct {
	mock.Mock
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
