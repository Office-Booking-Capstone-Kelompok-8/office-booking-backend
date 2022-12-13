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
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB, redisClient *redis.Client, conf *viper.Viper) {
	passwordService := passwordServicePkg.NewPasswordFuncImpl()
	validation := validator.NewValidator()
	generator := random.NewGenerator()
	imagekitService := imagekitServicePkg.NewImgKitService(conf.GetString("service.imgkit.privateKey"), conf.GetString("service.imgkit.publicKey"), conf.GetString("service.imgkit.endpoint"))
	mailService := mail.NewClient(conf.GetString("service.mailgun.domain"), conf.GetString("service.mailgun.apiKey"), conf.GetString("service.mailgun.sender"), conf.GetString("service.mailgun.senderName"))

	redisRepo := redisRepoPkg.NewRedisClient(redisClient)
	tokenService := authServicePkg.NewTokenServiceImpl(conf.GetString("token.access.secret"), conf.GetString("token.refresh.secret"), conf.GetDuration("token.access.exp"), conf.GetDuration("token.refresh.exp"), redisRepo)
	accessTokenMiddleware := middlewares.NewJWTMiddleware(conf.GetString("token.access.secret"), middlewares.ValidateAccessToken(tokenService))
	adminAccessTokenMiddleware := middlewares.NewJWTMiddleware(conf.GetString("token.access.secret"), middlewares.ValidateAdminAccessToken(tokenService))
	limiterMiddeleware := middlewares.NewLimiter(conf.GetDuration("otp.exp"))
	corsMiddleware := middlewares.NewCORSMiddleware(conf.GetStringSlice("server.allowedOrigins"))

	reservationRepository := reservationRepositoryPkg.NewReservationRepositoryImpl(db)
	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	authRepository := authRepositoryPkg.NewAuthRepositoryImpl(db)
	buildingRepository := buildingRepositoryPkg.NewBuildingRepositoryImpl(db)
	paymentRepository := paymentRepositoryPkg.NewPaymentRepositoryImpl(db)

	reservationService := reservationServicePkg.NewReservationServiceImpl(reservationRepository, buildingRepository)
	userService := userServicePkg.NewUserServiceImpl(userRepository, reservationService, imagekitService)
	buildingService := buildingServicePkg.NewBuildingServiceImpl(buildingRepository, reservationRepository, imagekitService, validation)
	authService := authServicePkg.NewAuthServiceImpl(authRepository, tokenService, redisRepo, mailService, passwordService, generator, conf)
	paymentService := paymentServicePkg.NewPaymentServiceImpl(paymentRepository, reservationRepository, imagekitService)

	reservationController := reservationControllerPkg.NewReservationController(reservationService, validation)
	userController := userControllerPkg.NewUserController(userService, validation)
	authController := authControllerPkg.NewAuthController(authService, validation)
	buildingController := buildingControllerPkg.NewBuildingController(buildingService, validation)
	paymentController := paymentControllerPkg.NewPaymentController(paymentService, validation)

	// init routes
	route := routes.NewRoutes(authController, userController, buildingController, reservationController, paymentController, limiterMiddeleware, accessTokenMiddleware, adminAccessTokenMiddleware, corsMiddleware)
	route.Init(app)
}
