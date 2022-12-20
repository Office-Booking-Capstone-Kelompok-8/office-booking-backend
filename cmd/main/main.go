package main

import (
	"context"
	"log"
	"net"
	"office-booking-backend/pkg/bootstrapper"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/custom"
	"office-booking-backend/pkg/database/mysql"
	"office-booking-backend/pkg/database/redis"
	"office-booking-backend/pkg/response"
	"office-booking-backend/pkg/utils/shutdown"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	local, err := time.LoadLocation(conf.GetString("server.timezone"))
	if err != nil {
		log.Fatal("failed to load timezone: ", err)
	}
	time.Local = local

	rs := redis.InitRedis(
		conf.GetString("service.redis.host"),
		conf.GetString("service.redis.port"),
		conf.GetString("service.redis.pass"),
		conf.GetString("service.redis.db"),
	)

	db := mysql.InitDatabase(
		conf.GetString("service.db.host"),
		conf.GetString("service.db.port"),
		conf.GetString("service.db.user"),
		conf.GetString("service.db.pass"),
		conf.GetString("service.db.name"),
	)

	app := fiber.New(fiber.Config{
		AppName:      conf.GetString("server.name"),
		ServerHeader: conf.GetString("server.header"),
		Prefork:      conf.GetBool("server.prefork"),
		ReadTimeout:  conf.GetDuration("server.read_timeout"),
		ErrorHandler: response.DefaultErrorHandler,
	})

	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType:        []fiber.ParserType{custom.CustomDate},
		ZeroEmpty:         true,
	})

	bootstrapper.InitAPI(app, db, rs, conf)

	wait := shutdown.GracefulShutdown(context.Background(), conf.GetDuration("server.shutdownTimeout"), map[string]shutdown.Operation{
		"fiber": func(ctx context.Context) error {
			return app.Shutdown()
		},

		"database": func(ctx context.Context) error {
			DB, err := db.DB()
			if err != nil {
				return err
			}
			return DB.Close()
		},

		"redis": func(ctx context.Context) error {
			return rs.Close()
		},
	})

	listen := net.JoinHostPort("", conf.GetString("server.port"))
	if err := app.Listen(listen); err != nil {
		log.Fatal(err)
	}

	<-wait
}
