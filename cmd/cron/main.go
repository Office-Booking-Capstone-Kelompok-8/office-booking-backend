package cron

import (
	"context"
	"log"
	"office-booking-backend/pkg/bootstrapper"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/database/mysql"
	"office-booking-backend/pkg/utils/shutdown"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	loc, err := time.LoadLocation(conf.GetString("server.timezone"))
	if err != nil {
		log.Fatal("failed to load timezone: ", err)
	}
	time.Local = loc

	db := mysql.InitDatabase(
		conf.GetString("service.db.host"),
		conf.GetString("service.db.port"),
		conf.GetString("service.db.user"),
		conf.GetString("service.db.pass"),
		conf.GetString("service.db.name"),
	)

	cron := gocron.NewScheduler(loc)

	bootstrapper.InitCron(db, cron, conf)

	wait := shutdown.GracefulShutdown(context.Background(), conf.GetDuration("server.shutdownTimeout"), map[string]shutdown.Operation{
		"database": func(ctx context.Context) error {
			DB, err := db.DB()
			if err != nil {
				return err
			}
			return DB.Close()
		},
		"cron": func(ctx context.Context) error {
			cron.Clear()
			return nil
		},
	})

	<-wait
}
