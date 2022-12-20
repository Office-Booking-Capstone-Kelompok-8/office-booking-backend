package impl

import (
	"context"
	"log"
	rr "office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/constant"
	"office-booking-backend/pkg/entity"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

type CronServiceImpl struct {
	reservation rr.ReservationRepository
	cron        *gocron.Scheduler
	conf        *viper.Viper
}

func NewCronServiceImpl(reservation rr.ReservationRepository, cron *gocron.Scheduler, conf *viper.Viper) *CronServiceImpl {
	return &CronServiceImpl{
		reservation: reservation,
		cron:        cron,
		conf:        conf,
	}
}

func (c *CronServiceImpl) Start() {
	c.cron.StartAsync()
	c.cron.Every(1).Day().At(c.conf.GetString("cron.executeAt")).Do(c.ScheduleReservationTask, context.Background())

	log.Println("cron service started")
}

func (c *CronServiceImpl) ScheduleReservationTask(ctx context.Context) error {
	reservations, err := c.reservation.GetReservationTaskUntilToday(ctx)
	if err != nil {
		log.Println("failed to get reservation task: ", err.Error())
		return err
	}

	for _, reservation := range *reservations {
		var err error
		if reservation.StatusID == constant.AWAITING_PAYMENT_STATUS {
			_, err = c.cron.At(reservation.ExpiredAt).Do(c.cancelReservation, ctx, reservation.ID)
		} else if reservation.StatusID == constant.ACTIVE_STATUS {
			_, err = c.cron.At(reservation.EndDate).Do(c.finishReservation, ctx, reservation.ID)
		}

		if err != nil {
			log.Println("failed to schedule reservation task: ", err.Error())
			return err
		}
	}

	return nil
}

func (c *CronServiceImpl) cancelReservation(ctx context.Context, reservationID string) error {
	return c.reservation.UpdateReservation(ctx, &entity.Reservation{
		ID:       reservationID,
		StatusID: constant.CANCELED_STATUS,
	})
}

func (c *CronServiceImpl) finishReservation(ctx context.Context, reservationID string) error {
	return c.reservation.UpdateReservation(ctx, &entity.Reservation{
		ID:       reservationID,
		StatusID: constant.COMPLETED_STATUS,
	})
}
