package mock

import (
	"office-booking-backend/internal/reservation/dto"

	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type ReservationServiceMock struct {
	mock.Mock
}

func (r *ReservationServiceMock) CountUserActiveReservations(ctx context.Context, userID string) (int64, error) {
	args := r.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ReservationServiceMock) CreateReservation(ctx context.Context, userID string, reservation *dto.AddReservartionRequest) error {
	args := r.Called(ctx, userID, reservation)
	return args.Error(0)
}
