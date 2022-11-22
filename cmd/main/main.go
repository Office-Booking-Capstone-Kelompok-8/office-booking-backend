package main

import (
	"context"
	"log"
	"office-booking-backend/pkg/bootstrapper"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/database"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("ENV") == "production" {
		return
	}

	//	load env variables from .env file for local development
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

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
	env := config.LoadConfig()

	redis := database.InitRedis(env["REDIS_HOST"], env["REDIS_PORT"], env["REDIS_PASS"], env["REDIS_DB"])
	db := database.InitDatabase(env["DB_HOST"], env["DB_PORT"], env["DB_USER"], env["DB_PASS"], env["DB_NAME"])

	app := fiber.New(fiber.Config{
		AppName:      config.APP_NAME,
		ServerHeader: config.SERVER_HEADER,
		Prefork:      env["PREFORK"] == "true",
		ReadTimeout:  time.Second * config.READ_TIMEOUT_SECONDS,
		ErrorHandler: config.DefaultErrorHandler,
	})

	bootstrapper.Init(app, db, redis)

	wait := gracefulShutdown(context.Background(), config.SHUTDOWN_TIMEOUT*time.Second, map[string]operation{
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
			return redis.Close()
		},
	})

	if err := app.Listen(":" + env["PORT"]); err != nil {
		log.Fatal(err)
	}

	<-wait
}
