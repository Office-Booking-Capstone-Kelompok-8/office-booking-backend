package bootstrapper

import (
	"office-booking-backend/pkg/routes"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, redisClient *redis.Client) {
	// TODO: Add your other initializations here

	// init routes
	route := routes.NewRoutes()
	route.Init(app)
}
