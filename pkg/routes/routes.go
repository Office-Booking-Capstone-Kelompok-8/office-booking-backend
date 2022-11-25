package routes

import (
	ac "office-booking-backend/internal/auth/controller"
	uc "office-booking-backend/internal/user/controller"
	"office-booking-backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	authController        *ac.AuthController
	userControllerPkg     *uc.UserController
	accessTokenMiddleware fiber.Handler
}

func NewRoutes(authController *ac.AuthController, userControllerPkg *uc.UserController, accessTokenMiddleware fiber.Handler) *Routes {
	return &Routes{
		authController:        authController,
		userControllerPkg:     userControllerPkg,
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

	user := v1.Group("/users")
	user.Get("/", r.accessTokenMiddleware, r.userControllerPkg.GetFullUserByID)
}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
