package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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

func Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
