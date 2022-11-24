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
	authController        *ac.AuthController
	accessTokenMiddleware fiber.Handler
}

func NewRoutes(authController *ac.AuthController, accessTokenMiddleware fiber.Handler) *Routes {
	return &Routes{
		authController:        authController,
		accessTokenMiddleware: accessTokenMiddleware,
	}
}

func (r *Routes) Init(app *fiber.App) {
	app.Use(recover.New(middlewares.RecoverConfig))
	app.Use(logger.New(middlewares.LoggerConfig))
	app.Use(cors.New(middlewares.CorsConfig))

	v1 := app.Group("/v1")
	v1.Get("/ping", ping)

	auth := v1.Group("/auth")
	auth.Post("/register", r.authController.RegisterUser)
	auth.Post("/login", r.authController.LoginUser)
	auth.Post("/logout", r.accessTokenMiddleware, r.authController.LogoutUser)
	auth.Post("/refresh", r.authController.RefreshToken)
}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
