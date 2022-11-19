package bootstrapper

import (
	"github.com/gofiber/fiber/v2"
	"office-booking-backend/pkg/routes"
)

func Init(app *fiber.App) {
	// TODO: Add your other initializations here

	// init routes
	route := routes.NewRoutes()
	route.Init(app)
}
