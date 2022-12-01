package routes

import (
	ac "office-booking-backend/internal/auth/controller"
	bc "office-booking-backend/internal/building/controller"
	uc "office-booking-backend/internal/user/controller"
	"office-booking-backend/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Routes struct {
	authController             *ac.AuthController
	userController             *uc.UserController
	buildingController         *bc.BuildingController
	accessTokenMiddleware      fiber.Handler
	adminAccessTokenMiddleware fiber.Handler
}

func NewRoutes(authController *ac.AuthController, userControllerPkg *uc.UserController, buildingController *bc.BuildingController, accessTokenMiddleware fiber.Handler, adminAccessTokenMiddleware fiber.Handler) *Routes {
	return &Routes{
		authController:             authController,
		userController:             userControllerPkg,
		buildingController:         buildingController,
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
	auth.Post("/request-otp", middlewares.OTPLimitter, r.authController.RequestOTP)
	auth.Post("/verify-otp", r.authController.VerifyOTP)

	// Enduser.User routes
	user := v1.Group("/users", r.accessTokenMiddleware)
	user.Get("/", r.userController.GetLoggedFullUserByID)
	user.Put("/", r.userController.UpdateLoggedUser)
	user.Put("/picture", r.userController.UpdateUserAvatar)
	user.Put("/change-password", r.authController.ChangePassword)

	// Admin routes
	admin := v1.Group("/admin", r.adminAccessTokenMiddleware)

	// Admin.User routes
	aUser := admin.Group("/users")
	aUser.Get("/", r.userController.GetAllUsers)
	aUser.Post("/new-admin", r.authController.RegisterAdmin)
	aUser.Post("/new-user", r.authController.RegisterUser)
	aUser.Get("/:userID", r.userController.GetFullUserByID)
	aUser.Put("/:userID", r.userController.UpdateUserByID)
	aUser.Delete("/:userID", r.userController.DeleteUserByID)

	// Admin.Building routes
	aBuilding := admin.Group("/buildings")
	aBuilding.Get("/", r.buildingController.GetAllBuildings)
	aBuilding.Get("/id", r.buildingController.RequestNewBuildingID)
	aBuilding.Get("/:buildingID", r.buildingController.GetBuildingDetailByID)
	aBuilding.Put("/:buildingID", r.buildingController.UpdateBuilding)
	aBuilding.Post("/:buildingID/picture", r.buildingController.UploadBuildingPicture)
	aBuilding.Post("/:buildingID/facilities", r.buildingController.AddBuildingFacilities)

	// Buildings routes
	building := v1.Group("/buildings")
	building.Get("/", r.buildingController.GetAllPublishedBuildings)
	building.Get("/:buildingID", r.buildingController.GetPublishedBuildingDetailByID)
	building.Get("/facilities/category", r.buildingController.GetFacilityCategories)

}

func ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
	})
}
