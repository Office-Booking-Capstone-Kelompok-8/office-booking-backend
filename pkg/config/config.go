package config

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	APP_NAME             = "office-zone-api v0.1"
	SERVER_HEADER        = "office-zone-api"
	READ_TIMEOUT_SECONDS = 10
	SHUTDOWN_TIMEOUT     = 15
)

var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Status code from errors if they implement *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Return status code with error JSON
	return c.Status(code).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func LoadConfig() map[string]string {
	env := make(map[string]string)

	env["DB_HOST"] = os.Getenv("DB_HOST")
	env["DB_PORT"] = os.Getenv("DB_PORT")
	env["DB_USER"] = os.Getenv("DB_USER")
	env["DB_PASS"] = os.Getenv("DB_PASS")
	env["DB_NAME"] = os.Getenv("DB_NAME")
	env["REDIS_HOST"] = os.Getenv("REDIS_HOST")
	env["REDIS_PORT"] = os.Getenv("REDIS_PORT")
	env["REDIS_PASSWORD"] = os.Getenv("REDIS_PASSWORD")
	env["REDIS_DB"] = os.Getenv("REDIS_DB")
	env["REFRESH_SECRET"] = os.Getenv("REFRESH_SECRET")
	env["ACCESS_SECRET"] = os.Getenv("ACCESS_SECRET")
	env["PORT"] = os.Getenv("PORT")
	env["PREFORK"] = os.Getenv("PREFORK")

	return env
}
