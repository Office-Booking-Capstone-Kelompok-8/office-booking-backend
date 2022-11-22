package bootstrapper

import (
	authControllerPkg "office-booking-backend/internal/auth/controller"
	authRepositoryPkg "office-booking-backend/internal/auth/repository/impl"
	authServicepkg "office-booking-backend/internal/auth/service/impl"
	"office-booking-backend/pkg/routes"
	passwordServicePkg "office-booking-backend/pkg/utils/password/impl"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB, redisClient *redis.Client) {
	passwordService := passwordServicePkg.NewPasswordFuncImpl()

	authRepository := authRepositoryPkg.NewAuthRepositoryImpl(db)
	authoService := authServicepkg.NewAuthServiceImpl(authRepository, passwordService)
	authController := authControllerPkg.NewAuthController(authoService)

	// init routes
	route := routes.NewRoutes(authController)
	route.Init(app)
}
