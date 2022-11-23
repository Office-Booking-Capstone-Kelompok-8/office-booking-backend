package routes

import (
	auth "office-booking-backend/internal/auth/controller"
	"office-booking-backend/pkg/handler"
	"office-booking-backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Routes struct {
	authController *auth.AuthController
}

func NewRoutes(authController *auth.AuthController) *Routes {
	return &Routes{
		authController: authController,
	}
}

func (r *Routes) Init(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New(middlewares.LoggerConfig))
	app.Use(cors.New(middlewares.CorsConfig))

	v1 := app.Group("/v1")
	v1.Get("/ping", handler.Ping)

	auth := v1.Group("/auth")
	auth.Post("/register", r.authController.RegisterUser)
	auth.Post("/login", r.authController.LoginUser)
}
