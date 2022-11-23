package routes

import (
	ac "office-booking-backend/internal/auth/controller"
	"office-booking-backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Routes struct {
	authController *ac.AuthController
}

func NewRoutes(authController *ac.AuthController) *Routes {
	return &Routes{
		authController: authController,
	}
}

func (r *Routes) Init(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New(middlewares.LoggerConfig))
	app.Use(cors.New(middlewares.CorsConfig))

	v1 := app.Group("/v1")
	v1.Get("/ping", ping)

	auth := v1.Group("/auth")
	auth.Post("/register", r.authController.RegisterUser)
	auth.Post("/login", r.authController.LoginUser)
}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
