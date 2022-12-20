package bootstrapper

import (
	cronServicePkg "office-booking-backend/internal/cron/service/impl"
	reservationRepositoryPkg "office-booking-backend/internal/reservation/repository/impl"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitCron(db *gorm.DB, cron *gocron.Scheduler, conf *viper.Viper) {
	reservationRepository := reservationRepositoryPkg.NewReservationRepositoryImpl(db)
	cronService := cronServicePkg.NewCronServiceImpl(reservationRepository, cron, conf)
	cronService.Start()
}
