package service

import "context"

type CronService interface {
	ScheduleReservationTask(ctx context.Context) error
	Start()
}
