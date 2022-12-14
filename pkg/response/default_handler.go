package response

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var DefaultErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Status code from errors if they implement *fiber.Error
	var e *fiber.Error
	message := err.Error()
	if errors.As(err, &e) {
		if strings.Contains(message, "Cannot") && e.Code == fiber.StatusNotFound {
			message = "Not Found"
		}

		code = e.Code
	}

	// Return status code with error JSON
	return c.Status(code).JSON(BaseResponse{
		Message: message,
	})
}
