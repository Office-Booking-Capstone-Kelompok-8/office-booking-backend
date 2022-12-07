package bootstrapper

import (
	authControllerPkg "office-booking-backend/internal/auth/controller"
	authRepositoryPkg "office-booking-backend/internal/auth/repository/impl"
	authServicePkg "office-booking-backend/internal/auth/service/impl"
	buildingControllerPkg "office-booking-backend/internal/building/controller"
	buildingRepositoryPkg "office-booking-backend/internal/building/repository/impl"
	buildingServicePkg "office-booking-backend/internal/building/service/impl"
	paymentControllerPkg "office-booking-backend/internal/payment/controller"
	paymentRepositoryPkg "office-booking-backend/internal/payment/repository/impl"
	paymentServicePkg "office-booking-backend/internal/payment/service/impl"
	reservationControllerPkg "office-booking-backend/internal/reservation/controller"
	reservationRepositoryPkg "office-booking-backend/internal/reservation/repository/impl"
	reservationServicePkg "office-booking-backend/internal/reservation/service/impl"
	userControllerPkg "office-booking-backend/internal/user/controller"
	userRepositoryPkg "office-booking-backend/internal/user/repository/impl"
	userServicePkg "office-booking-backend/internal/user/service/impl"

	"office-booking-backend/pkg/config"
	redisRepoPkg "office-booking-backend/pkg/database/redis"
	"office-booking-backend/pkg/middlewares"
	"office-booking-backend/pkg/routes"
	imagekitServicePkg "office-booking-backend/pkg/utils/imagekit"
	"office-booking-backend/pkg/utils/mail"
	passwordServicePkg "office-booking-backend/pkg/utils/password"
	"office-booking-backend/pkg/utils/random"
	"office-booking-backend/pkg/utils/validator"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB, redisClient *redis.Client, conf map[string]string) {

	imagekitService := imagekitServicePkg.NewImgKitService(conf["IMGKIT_PRIVATE_KEY"], conf["IMGKIT_PUBLIC_KEY"], conf["IMGKIT_URL_ENDPOINT"])
	passwordService := passwordServicePkg.NewPasswordFuncImpl()
	validation := validator.NewValidator()
	generator := random.NewGenerator()
	redisRepo := redisRepoPkg.NewRedisClient(redisClient)
	mailService := mail.NewClient(conf["MAIL_DOMAIN"], conf["MAIL_API_KEY"], conf["MAIL_SENDER"], conf["MAIL_SENDER_NAME"])
	tokenService := authServicePkg.NewTokenServiceImpl(conf["ACCESS_SECRET"], conf["REFRESH_SECRET"], config.ACCESS_TOKEN_DURATION, config.REFRESH_TOKEN_DURATION, redisRepo)

	accessTokenMiddleware := middlewares.NewJWTMiddleware(conf["ACCESS_SECRET"], middlewares.ValidateAccessToken(tokenService))
	adminAccessTokenMiddleware := middlewares.NewJWTMiddleware(conf["ACCESS_SECRET"], middlewares.ValidateAdminAccessToken(tokenService))

	reservationRepository := reservationRepositoryPkg.NewReservationRepositoryImpl(db)
	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	authRepository := authRepositoryPkg.NewAuthRepositoryImpl(db)
	buildingRepository := buildingRepositoryPkg.NewBuildingRepositoryImpl(db)
	paymentRepository := paymentRepositoryPkg.NewPaymentRepositoryImpl(db)

	reservationService := reservationServicePkg.NewReservationServiceImpl(reservationRepository, buildingRepository)
	userService := userServicePkg.NewUserServiceImpl(userRepository, reservationService, imagekitService)
	buildingService := buildingServicePkg.NewBuildingServiceImpl(buildingRepository, reservationRepository, imagekitService, validation)
	authService := authServicePkg.NewAuthServiceImpl(authRepository, userRepository, tokenService, redisRepo, mailService, passwordService, generator)
	paymentService := paymentServicePkg.NewPaymentServiceImpl(paymentRepository)

	reservationController := reservationControllerPkg.NewReservationController(reservationService, validation)
	userController := userControllerPkg.NewUserController(userService, validation)
	authController := authControllerPkg.NewAuthController(authService, validation)
	buildingController := buildingControllerPkg.NewBuildingController(buildingService, validation)
	paymentController := paymentControllerPkg.NewPaymentController(paymentService, validation)

	// init routes
	route := routes.NewRoutes(authController, userController, buildingController, reservationController, paymentController, accessTokenMiddleware, adminAccessTokenMiddleware)
	route.Init(app)
}
