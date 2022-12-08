package routes

import (
	ac "office-booking-backend/internal/auth/controller"
	bc "office-booking-backend/internal/building/controller"
	pr "office-booking-backend/internal/payment/controller"
	rc "office-booking-backend/internal/reservation/controller"
	uc "office-booking-backend/internal/user/controller"
	"office-booking-backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	authController             *ac.AuthController
	userController             *uc.UserController
	buildingController         *bc.BuildingController
	reservationController      *rc.ReservationController
	paymentController          *pr.PaymentController
	accessTokenMiddleware      fiber.Handler
	adminAccessTokenMiddleware fiber.Handler
}

func NewRoutes(authController *ac.AuthController, userControllerPkg *uc.UserController, buildingController *bc.BuildingController, reservationController *rc.ReservationController, paymentController *pr.PaymentController, accessTokenMiddleware fiber.Handler, adminAccessTokenMiddleware fiber.Handler) *Routes {
	return &Routes{
		authController:             authController,
		userController:             userControllerPkg,
		buildingController:         buildingController,
		reservationController:      reservationController,
		paymentController:          paymentController,
		accessTokenMiddleware:      accessTokenMiddleware,
		adminAccessTokenMiddleware: adminAccessTokenMiddleware,
	}
}

func (r *Routes) Init(app *fiber.App) {
	app.Use(middlewares.Recover)
	app.Use(middlewares.Logger)
	app.Use(middlewares.Cors)

	v1 := app.Group("/v1")
	v1.Get("/ping", ping)

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", r.authController.RegisterUser)
	auth.Post("/login", r.authController.LoginUser)
	auth.Post("/logout", r.accessTokenMiddleware, r.authController.LogoutUser)
	auth.Post("/refresh", r.authController.RefreshToken)
	auth.Put("/reset-password", r.authController.ResetPassword)

	otp := auth.Group("/otp")
	otp.Post("/request", middlewares.OTPLimitter, r.authController.RequestOTP)
	otp.Post("/verify", r.authController.VerifyOTP)

	// Enduser.User routes
	user := v1.Group("/users", r.accessTokenMiddleware)
	user.Get("/", r.userController.GetLoggedFullUserByID)
	user.Put("/", r.userController.UpdateLoggedUser)
	user.Put("/picture", r.userController.UpdateUserAvatar)
	user.Put("/change-password", r.authController.ChangePassword)

	// Enduser.Reservation routes
	reservation := v1.Group("/reservations", r.accessTokenMiddleware)
	reservation.Get("/", r.reservationController.GetUserReservations)
	reservation.Post("/", r.reservationController.CreateReservation)
	reservation.Get("/:reservationID", r.reservationController.GetUserReservationDetailByID)
	reservation.Delete("/:reservationID", r.reservationController.CancelReservation)

	// Admin routes
	admin := v1.Group("/admin", r.adminAccessTokenMiddleware)

	// Admin.User routes
	aUser := admin.Group("/users")
	aUser.Get("/", r.userController.GetAllUsers)
	aUser.Post("/", r.authController.RegisterUser)
	aUser.Post("/admin", r.authController.RegisterAdmin)
	aUser.Get("/:userID", r.userController.GetFullUserByID)
	aUser.Put("/:userID", r.userController.UpdateUserByID)
	aUser.Delete("/:userID", r.userController.DeleteUserByID)
	aUser.Put("/:userID/picture", r.userController.UpdateAnotherUserAvatar)

	// Admin.Building routes
	aBuilding := admin.Group("/buildings")
	aBuilding.Get("/", r.buildingController.GetAllBuildings)
	aBuilding.Get("/id", r.buildingController.RequestNewBuildingID)
	aBuilding.Get("/:buildingID", r.buildingController.GetBuildingDetailByID)
	aBuilding.Put("/:buildingID", r.buildingController.UpdateBuilding)
	aBuilding.Delete("/:buildingID", r.buildingController.DeleteBuilding)
	aBuilding.Post("/:buildingID/pictures", r.buildingController.AddBuildingPicture)
	aBuilding.Delete("/:buildingID/pictures/:pictureID", r.buildingController.DeleteBuildingPicture)
	aBuilding.Post("/:buildingID/facilities", r.buildingController.AddBuildingFacilities)
	aBuilding.Delete("/:buildingID/facilities/:facilityID", r.buildingController.DeleteBuildingFacility)

	// Admin.Reservation routes
	aReservation := admin.Group("/reservations")
	aReservation.Get("/", r.reservationController.GetReservations)
	aReservation.Post("/", r.reservationController.CreateAdminReservation)
	aReservation.Get("/:reservationID", r.reservationController.GetReservationDetailByID)
	aReservation.Put("/:reservationID", r.reservationController.UpdateReservation)
	aReservation.Delete("/:reservationID", r.reservationController.DeleteReservation)

	// Admin.Payment routes
	aPayment := admin.Group("/payments", r.adminAccessTokenMiddleware)
	aPayment.Post("/", r.paymentController.CreatePayment)
	aPayment.Get("/:paymentID", r.paymentController.GetPaymentByID)
	aPayment.Put("/:paymentID", r.paymentController.UpdatePayment)

	// Buildings routes
	building := v1.Group("/buildings")
	building.Get("/", r.buildingController.GetAllPublishedBuildings)
	building.Get("/facilities/category", r.buildingController.GetFacilityCategories)
	building.Get("/:buildingID", r.buildingController.GetPublishedBuildingDetailByID)

	// Location routes
	location := v1.Group("/locations")
	location.Get("/cities", r.buildingController.GetCities)
	location.Get("/districts", r.buildingController.GetDistricts)

	// Payment routes
	payment := v1.Group("/payments")
	payment.Get("/banks", r.paymentController.GetBanks)
	payment.Get("/:paymentID", r.paymentController.GetPaymentByID)
}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
