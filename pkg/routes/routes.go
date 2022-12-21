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
	auth                       *ac.AuthController
	user                       *uc.UserController
	building                   *bc.BuildingController
	reservation                *rc.ReservationController
	payment                    *pr.PaymentController
	limiter                    *middlewares.Limiter
	cors                       fiber.Handler
	accessTokenMiddleware      fiber.Handler
	adminAccessTokenMiddleware fiber.Handler
}

func NewRoutes(authController *ac.AuthController, userControllerPkg *uc.UserController, buildingController *bc.BuildingController, reservationController *rc.ReservationController, paymentController *pr.PaymentController, limiter *middlewares.Limiter, accessTokenMiddleware fiber.Handler, adminAccessTokenMiddleware fiber.Handler, cors fiber.Handler) *Routes {
	return &Routes{
		auth:                       authController,
		user:                       userControllerPkg,
		building:                   buildingController,
		reservation:                reservationController,
		payment:                    paymentController,
		limiter:                    limiter,
		cors:                       cors,
		accessTokenMiddleware:      accessTokenMiddleware,
		adminAccessTokenMiddleware: adminAccessTokenMiddleware,
	}
}

func (r *Routes) Init(app *fiber.App) {
	app.Use(middlewares.Recover)
	app.Use(middlewares.Logger)
	app.Use(r.cors)

	v1 := app.Group("/v1")
	v1.Get("/ping", ping)

	// Payment routes
	payment := v1.Group("/payments")
	payment.Get("/methods", r.payment.GetAllPaymentMethod)
	payment.Get("/methods/banks", r.payment.GetBanks)
	payment.Get("/methods/:paymentMethodID", r.payment.GetPaymentMethodByID)

	// Buildings routes
	building := v1.Group("/buildings")
	building.Get("/", r.building.GetAllPublishedBuildings)
	building.Get("/facilities/category", r.building.GetFacilityCategories)
	building.Get("/:buildingID", r.building.GetPublishedBuildingDetailByID)
	building.Get("/:buildingID/reviews", r.building.GetBuildingReviews)

	// Location routes
	location := v1.Group("/locations")
	location.Get("/cities", r.building.GetCities)
	location.Get("/districts", r.building.GetDistricts)

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", r.auth.RegisterUser)
	auth.Post("/login", r.auth.LoginUser)
	auth.Post("/logout", r.accessTokenMiddleware, r.auth.LogoutUser)
	auth.Post("/refresh", r.auth.RefreshToken)
	auth.Put("/reset-password", r.auth.ResetPassword)

	otp := auth.Group("/otp")
	requestOtp := otp.Group("/request")
	requestOtp.Post("/reset-password", r.limiter.ResetPasswordOTPLimitter(), r.auth.RequestPasswordResetOTP)
	requestOtp.Post("/email", r.accessTokenMiddleware, r.limiter.ResetPasswordOTPLimitter(), r.auth.RequestVerifyEmailOTP)
	verifyOtp := otp.Group("/verify")
	verifyOtp.Post("/reset-password", r.auth.VerifyPasswordResetOTP)
	verifyOtp.Post("/email", r.accessTokenMiddleware, r.auth.VerifyEmailOTP)

	// Enduser.User routes
	user := v1.Group("/users")
	user.Get("/", r.accessTokenMiddleware, r.user.GetLoggedFullUserByID)
	user.Put("/", r.accessTokenMiddleware, r.user.UpdateLoggedUser)
	user.Put("/picture", r.accessTokenMiddleware, r.user.UpdateUserAvatar)
	user.Put("/change-password", r.accessTokenMiddleware, r.auth.ChangePassword)

	// Enduser.Reservation routes
	uReservation := v1.Group("/reservations")
	uReservation.Get("/", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.GetUserReservations)
	uReservation.Post("/", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.CreateReservation)
	uReservation.Get("/:reservationID", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.GetUserReservationDetailByID)
	uReservation.Delete("/:reservationID", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.CancelReservation)
	uReservation.Post("/:reservationID/reviews", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.CreateReservationReview)
	uReservation.Put("/:reservationID/reviews", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.UpdateReservationReview)
	uReservation.Get("/:reservationID/reviews", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.reservation.GetUserReservationReview)

	// Enduser.Payment routes
	payment.Get("/:reservationID", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.payment.GetUsereservationPaymentByID)
	payment.Post("/:reservationID", r.accessTokenMiddleware, middlewares.EnforceValidEmail(), r.payment.UploadPaymentProof)

	// Admin routes
	admin := v1.Group("/admin")

	// Admin.User routes
	aUser := admin.Group("/users")
	aUser.Get("/", r.adminAccessTokenMiddleware, r.user.GetAllUsers)
	aUser.Post("/", r.adminAccessTokenMiddleware, r.auth.RegisterUser)
	aUser.Post("/admin", r.adminAccessTokenMiddleware, r.auth.RegisterAdmin)
	aUser.Get("/total", r.adminAccessTokenMiddleware, r.user.GetRegisteredMemberCount)
	aUser.Get("/statistics", r.adminAccessTokenMiddleware, r.user.GetRegisteredMemberStat)
	aUser.Get("/:userID", r.adminAccessTokenMiddleware, r.user.GetFullUserByID)
	aUser.Put("/:userID", r.adminAccessTokenMiddleware, r.user.UpdateUserByID)
	aUser.Delete("/:userID", r.adminAccessTokenMiddleware, r.user.DeleteUserByID)
	aUser.Put("/:userID/picture", r.adminAccessTokenMiddleware, r.user.UpdateAnotherUserAvatar)

	// Admin.Building routes
	aBuilding := admin.Group("/buildings")
	aBuilding.Get("/", r.adminAccessTokenMiddleware, r.building.GetAllBuildings)
	aBuilding.Get("/id", r.adminAccessTokenMiddleware, r.building.RequestNewBuildingID)
	aBuilding.Get("/total", r.adminAccessTokenMiddleware, r.building.GetBuildingTotal)
	aBuilding.Get("/:buildingID", r.adminAccessTokenMiddleware, r.building.GetBuildingDetailByID)
	aBuilding.Delete("/:buildingID", r.adminAccessTokenMiddleware, r.building.DeleteBuilding)
	aBuilding.Put("/:buildingID", r.adminAccessTokenMiddleware, r.building.UpdateBuilding)
	aBuilding.Put("/:buildingID/publish", r.adminAccessTokenMiddleware, r.building.UpdateBuildingPublishState)
	aBuilding.Post("/:buildingID/facilities", r.adminAccessTokenMiddleware, r.building.AddBuildingFacilities)
	aBuilding.Delete("/:buildingID/facilities/:facilityID", r.adminAccessTokenMiddleware, r.building.DeleteBuildingFacility)
	aBuilding.Post("/:buildingID/pictures", r.adminAccessTokenMiddleware, r.building.AddBuildingPicture)
	aBuilding.Delete("/:buildingID/pictures/:pictureID", r.adminAccessTokenMiddleware, r.building.DeleteBuildingPicture)

	// Admin.Reservation routes
	aReservation := admin.Group("/reservations")
	aReservation.Get("/", r.adminAccessTokenMiddleware, r.reservation.GetReservations)
	aReservation.Post("/", r.adminAccessTokenMiddleware, r.reservation.CreateAdminReservation)
	aReservation.Get("/total", r.adminAccessTokenMiddleware, r.reservation.GetReservationTotal)
	aReservation.Get("/revenue", r.adminAccessTokenMiddleware, r.reservation.GetTotalRevenueByTime)
	aReservation.Get("/:reservationID", r.adminAccessTokenMiddleware, r.reservation.GetReservationDetailByID)
	aReservation.Put("/:reservationID", r.adminAccessTokenMiddleware, r.reservation.UpdateReservation)
	aReservation.Delete("/:reservationID", r.adminAccessTokenMiddleware, r.reservation.DeleteReservation)
	aReservation.Put("/:reservationID/status", r.adminAccessTokenMiddleware, r.reservation.UpdateReservationStatus)

	// Admin.Payment routes
	aPayment := admin.Group("/payments")
	aPayment.Post("/", r.adminAccessTokenMiddleware, r.payment.CreatePaymentMethod)
	aPayment.Put("/:paymentID", r.adminAccessTokenMiddleware, r.payment.UpdatePaymentMethod)
	aPayment.Delete("/:paymentID", r.adminAccessTokenMiddleware, r.payment.DeletePaymentMethod)
	aPayment.Get("/reservations/:reservationID/", r.adminAccessTokenMiddleware, r.payment.GetReservationPaymentByID)
}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
