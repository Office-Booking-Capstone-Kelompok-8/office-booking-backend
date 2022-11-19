package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"office-booking-backend/pkg/bootstrapper"
	"office-booking-backend/pkg/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	app := fiber.New(fiber.Config{
		AppName:      config.APP_NAME,
		ServerHeader: config.SERVER_HEADER,
		Prefork:      env["PREFORK"] == "true",
		ReadTimeout:  time.Second * config.READ_TIMEOUT_SECONDS,
		ErrorHandler: config.DefaultErrorHandler,
	})

	bootstrapper.Init(app)

	wait := gracefulShutdown(context.Background(), config.SHUTDOWN_TIMEOUT*time.Second, map[string]operation{
		// Add cleanup operations here

		// TODO: add database connection clean up
		// "database": func(ctx context.Context) error {
		// 	return db.Close()
		// },

		"fiber": func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	if err := app.Listen(":" + env["PORT"]); err != nil {
		log.Fatal(err)
	}

	<-wait
}
