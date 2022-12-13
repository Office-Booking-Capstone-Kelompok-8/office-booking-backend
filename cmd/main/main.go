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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

// operation is a cleanup function on shutting down
type operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
// credit: Alfian Dhimas
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	bootstrapper.Init(app, db, rs, conf)

	wait := gracefulShutdown(context.Background(), conf.GetDuration("server.shutdownTimeout"), map[string]operation{
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
