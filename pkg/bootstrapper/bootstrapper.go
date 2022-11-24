package bootstrapper

import (
	authControllerPkg "office-booking-backend/internal/auth/controller"
	authRepositoryPkg "office-booking-backend/internal/auth/repository/impl"
	authServicepkg "office-booking-backend/internal/auth/service/impl"
	"office-booking-backend/pkg/config"
	"office-booking-backend/pkg/middlewares"
	"office-booking-backend/pkg/routes"
	passwordServicePkg "office-booking-backend/pkg/utils/password/impl"
	"office-booking-backend/pkg/utils/validator"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB, redisClient *redis.Client, conf map[string]string) {
	passwordService := passwordServicePkg.NewPasswordFuncImpl()

	validation := validator.NewValidator()

	authRepository := authRepositoryPkg.NewAuthRepositoryImpl(db)
	tokenRepository := authRepositoryPkg.NewTokenRepositoryImpl(redisClient)
	tokenService := authServicepkg.NewTokenServiceImpl(conf["ACCESS_SECRET"], conf["REFRESH_SECRET"], config.ACCESS_TOKEN_DURATION, config.REFRESH_TOKEN_DURATION, tokenRepository)
	accessTokenMiddleware := middlewares.NewJWTMiddleware(conf["ACCESS_SECRET"], middlewares.ValidateAccessToken(tokenService))

	authService := authServicepkg.NewAuthServiceImpl(authRepository, tokenService, passwordService)
	authController := authControllerPkg.NewAuthController(authService, validation)

	// init routes
	route := routes.NewRoutes(authController, accessTokenMiddleware)
	route.Init(app)
}
