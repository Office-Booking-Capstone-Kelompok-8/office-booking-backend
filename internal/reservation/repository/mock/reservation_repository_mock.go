package mock

import (
	"office-booking-backend/internal/reservation/dto"
	"office-booking-backend/pkg/entity"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
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

func (r *ReservationRepositoryMock) CountReservation(ctx context.Context, filter *dto.ReservationQueryParam) (int64, error) {
	args := r.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ReservationRepositoryMock) IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time, excludedReservationID ...string) (bool, error) {
	args := r.Called(ctx, buildingID, start, end, excludedReservationID)
	return args.Get(0).(bool), args.Error(1)
}

func (r *ReservationRepositoryMock) GetUserReservationByID(ctx context.Context, reservationID string, userID string) (*entity.Reservation, error) {
	args := r.Called(ctx, reservationID, userID)
	return args.Get(0).(*entity.Reservation), args.Error(1)
}

func (r *ReservationRepositoryMock) GetReservationCountByStatus(ctx context.Context) (*entity.StatusesStat, error) {
	args := r.Called(ctx)
	return args.Get(0).(*entity.StatusesStat), args.Error(1)
}

func (r *ReservationRepositoryMock) GetReservationCountByTime(ctx context.Context) (*entity.TimeframeStat, error) {
	args := r.Called(ctx)
	return args.Get(0).(*entity.TimeframeStat), args.Error(1)
}

func (r *ReservationRepositoryMock) GetTotalRevenue(ctx context.Context) (*entity.TimeframeStat, error) {
	args := r.Called(ctx)
	return args.Get(0).(*entity.TimeframeStat), args.Error(1)
}

func (r *ReservationRepositoryMock) GetReservationTaskUntilToday(ctx context.Context) (*entity.Reservations, error) {
	args := r.Called(ctx)
	return args.Get(0).(*entity.Reservations), args.Error(1)
}

func (r *ReservationRepositoryMock) DeleteReservationByID(ctx context.Context, reservationID string) error {
	args := r.Called(ctx, reservationID)
	return args.Error(0)
}

func (r *ReservationRepositoryMock) UpdateReservation(ctx context.Context, reservation *entity.Reservation) error {
	args := r.Called(ctx, reservation)
	return args.Error(0)
}

func (r *ReservationRepositoryMock) GetReservations(ctx context.Context, filter *dto.ReservationQueryParam) (*entity.Reservations, error) {
	args := r.Called(ctx, filter)
	return args.Get(0).(*entity.Reservations), args.Error(1)
}
