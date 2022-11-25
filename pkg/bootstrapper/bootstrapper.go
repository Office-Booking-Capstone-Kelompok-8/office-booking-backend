package bootstrapper

import (
	authControllerPkg "office-booking-backend/internal/auth/controller"
	authRepositoryPkg "office-booking-backend/internal/auth/repository/impl"
	authServicePkg "office-booking-backend/internal/auth/service/impl"
	userControllerPkg "office-booking-backend/internal/user/controller"
	userRepositoryPkg "office-booking-backend/internal/user/repository/impl"
	userServicePkg "office-booking-backend/internal/user/service/impl"
	"office-booking-backend/pkg/config"
	redisRepoPkg "office-booking-backend/pkg/database/redis"
	"office-booking-backend/pkg/middlewares"
	"office-booking-backend/pkg/routes"
	"office-booking-backend/pkg/utils/mail"
	passwordServicePkg "office-booking-backend/pkg/utils/password"
	"office-booking-backend/pkg/utils/random"
	"office-booking-backend/pkg/utils/validator"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB, redisClient *redis.Client, conf map[string]string) {
	passwordService := passwordServicePkg.NewPasswordFuncImpl()

	validation := validator.NewValidator()
	generator := random.NewGenerator()
	redisRepo := redisRepoPkg.NewRedisClient(redisClient)
	mailService := mail.NewClient(conf["MAIL_DOMAIN"], conf["MAIL_API_KEY"], conf["MAIL_SENDER"], conf["MAIL_SENDER_NAME"])
	tokenService := authServicePkg.NewTokenServiceImpl(conf["ACCESS_SECRET"], conf["REFRESH_SECRET"], config.ACCESS_TOKEN_DURATION, config.REFRESH_TOKEN_DURATION, redisRepo)

	accessTokenMiddleware := middlewares.NewJWTMiddleware(conf["ACCESS_SECRET"], middlewares.ValidateAccessToken(tokenService))
	adminAccessTokenMiddleware := middlewares.NewJWTMiddleware(conf["ACCESS_SECRET"], middlewares.ValidateAdminAccessToken(tokenService))

	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository)
	userController := userControllerPkg.NewUserController(userService)

	authRepository := authRepositoryPkg.NewAuthRepositoryImpl(db)
	authService := authServicePkg.NewAuthServiceImpl(authRepository, userRepository, tokenService, redisRepo, mailService, passwordService, generator)
	authController := authControllerPkg.NewAuthController(authService, validation)

	// init routes
	route := routes.NewRoutes(authController, userController, accessTokenMiddleware, adminAccessTokenMiddleware)
	route.Init(app)
}
