package middlewares

import "github.com/gofiber/fiber/v2/middleware/cors"

var CorsConfig = cors.Config{
	AllowOrigins: "*",
	AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
}
