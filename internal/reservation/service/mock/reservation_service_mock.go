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

func (r *ReservationServiceMock) IsBuildingAvailable(ctx context.Context, buildingID string, start time.Time, end time.Time) (bool, error) {
	args := r.Called(ctx, buildingID, start, end)
	return args.Get(0).(bool), args.Error(1)
}

func (r *ReservationServiceMock) CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) (string, error) {
	args := r.Called(ctx, userID, reservation)
	return args.Get(0).(string), args.Error(1)
}
