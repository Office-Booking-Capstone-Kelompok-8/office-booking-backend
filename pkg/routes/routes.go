package routes

import (
	ac "office-booking-backend/internal/auth/controller"
	"office-booking-backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
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
	app.Use(middlewares.Recover)
	app.Use(middlewares.Logger)
	app.Use(middlewares.Cors)

	v1 := app.Group("/v1")
	v1.Get("/ping", ping)

	auth := v1.Group("/auth")
	auth.Post("/register", r.authController.RegisterUser)
	auth.Post("/login", r.authController.LoginUser)
	auth.Post("/logout", r.accessTokenMiddleware, r.authController.LogoutUser)
	auth.Post("/refresh", r.authController.RefreshToken)
	auth.Put("/reset-password", r.authController.ResetPassword)
	auth.Post("/request-otp", middlewares.OTPLimitter, r.authController.RequestOTP)
	auth.Post("/verify-otp", r.authController.VerifyOTP)
}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
