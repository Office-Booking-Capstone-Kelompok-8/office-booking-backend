package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewCORSMiddleware(allowedOrigins []string) fiber.Handler {
	allowed := ""
	for _, origin := range allowedOrigins {
		allowed += fmt.Sprintf("%s,", origin)
	}

	return cors.New(cors.Config{
		AllowOrigins: allowed,
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
