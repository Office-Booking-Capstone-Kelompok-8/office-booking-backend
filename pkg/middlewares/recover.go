package middlewares

import (
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var Recover = recover.New(recover.Config{
	EnableStackTrace: true,
})
