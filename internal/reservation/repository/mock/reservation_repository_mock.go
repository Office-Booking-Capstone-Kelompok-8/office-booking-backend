package mock

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"office-booking-backend/pkg/entity"
	"time"
)

type ReservationRepositoryMock struct {
	mock.Mock
}

func (r *ReservationRepositoryMock) CountUserActiveReservations(ctx context.Context, userID string) (int64, error) {
	args := r.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ReservationRepositoryMock) CountBuildingActiveReservations(ctx context.Context, buildingID string) (int64, error) {
	args := r.Called(ctx, buildingID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ReservationRepositoryMock) IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error) {
	args := r.Called(ctx, buildingID, start, end)
	return args.Get(0).(bool), args.Error(1)
}

func (r *ReservationRepositoryMock) AddBuildingReservation(ctx context.Context, reservation *entity.Reservation) error {
	args := r.Called(ctx, reservation)
	return args.Error(0)
}

func (r *ReservationRepositoryMock) CountUserReservation(ctx context.Context, userID string) (int64, error) {
	args := r.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ReservationRepositoryMock) GetUserReservations(ctx context.Context, userID string, offset int, limit int) (*entity.Reservations, error) {
	args := r.Called(ctx, userID, offset, limit)
	return args.Get(0).(*entity.Reservations), args.Error(1)
}

func (r *ReservationRepositoryMock) GetReservationByID(ctx context.Context, reservationID string) (*entity.Reservation, error) {
	args := r.Called(ctx, reservationID)
	return args.Get(0).(*entity.Reservation), args.Error(1)
}

func (r *ReservationRepositoryMock) DeleteReservationByID(ctx context.Context, reservationID string) error {
	args := r.Called(ctx, reservationID)
	return args.Error(0)
}

func (r *ReservationRepositoryMock) UpdateReservation(ctx context.Context, reservation *entity.Reservation) error {
	args := r.Called(ctx, reservation)
	return args.Error(0)
}
