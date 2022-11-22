package bootstrapper

import (
	"office-booking-backend/pkg/routes"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	// TODO: Add your other initializations here

	// init routes
	route := routes.NewRoutes()
	route.Init(app)
}
