package impl

import (
	"context"
	"log"
	pr "office-booking-backend/internal/payment/repository"
	rr "office-booking-backend/internal/reservation/repository"
	"office-booking-backend/pkg/constant"
	"office-booking-backend/pkg/entity"
	"time"

	err2 "office-booking-backend/pkg/errors"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

type CronServiceImpl struct {
	reservation rr.ReservationRepository
	payment     pr.PaymentRepository
	cron        *gocron.Scheduler
	conf        *viper.Viper
}

func NewCronServiceImpl(reservation rr.ReservationRepository, payment pr.PaymentRepository, cron *gocron.Scheduler, conf *viper.Viper) *CronServiceImpl {
	return &CronServiceImpl{
		reservation: reservation,
		payment:     payment,
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
			c.scheduleCancelReservation(ctx, reservation.ID, reservation.ExpiredAt)
		} else if reservation.StatusID == constant.ACTIVE_STATUS {
			c.scheduleFinishReservation(ctx, reservation.ID, reservation.EndDate)
		}

		if err != nil {
			log.Println("failed to schedule reservation task: ", err.Error())
			return err
		}
	}

	return nil
}

func (c *CronServiceImpl) scheduleCancelReservation(ctx context.Context, reservationID string, executeAt time.Time) error {
	if executeAt.Before(time.Now()) {
		return c.cancelReservation(ctx, reservationID)
	}

	_, err := c.cron.At(executeAt).Do(c.cancelReservation, ctx, reservationID)
	return err
}

func (c *CronServiceImpl) scheduleFinishReservation(ctx context.Context, reservationID string, executeAt time.Time) error {
	if executeAt.Before(time.Now()) {
		return c.finishReservation(ctx, reservationID)
	}

	_, err := c.cron.At(executeAt).Do(c.cancelReservation, ctx, reservationID)
	return err
}

func (c *CronServiceImpl) cancelReservation(ctx context.Context, reservationID string) error {
	_, err := c.payment.GetReservationPaymentByID(ctx, reservationID, "")
	if err == err2.ErrPaymentNotFound {
		return c.reservation.UpdateReservation(ctx, &entity.Reservation{
			ID:       reservationID,
			StatusID: constant.CANCELED_STATUS,
		})
	}

	return err
}

func (c *CronServiceImpl) finishReservation(ctx context.Context, reservationID string) error {
	return c.reservation.UpdateReservation(ctx, &entity.Reservation{
		ID:       reservationID,
		StatusID: constant.COMPLETED_STATUS,
	})
}
