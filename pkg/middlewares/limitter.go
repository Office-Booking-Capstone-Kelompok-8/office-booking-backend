package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/pkg/config"
)

var ResetPasswordOTPLimitter = limiter.New(limiter.Config{
	Max:        1,
	Expiration: config.OTP_RESEND_TIME,
	KeyGenerator: func(c *fiber.Ctx) string {
		// limit by ip
		// ip := c.Get("x-real-ip")
		// if ip == "" {
		//  	return c.IP()
		// }
		// return ip

		// Limit with email
		body := new(struct {
			Email string `json:"email"`
		})
		if err := c.BodyParser(body); err != nil {
			return ""
		}

		return body.Email
	},
	LimitReached: func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTooManyRequests, "too many requests")
	},
})

var VerifyEmailOTPLimitter = limiter.New(limiter.Config{
	Max:        1,
	Expiration: config.OTP_RESEND_TIME,
	KeyGenerator: func(c *fiber.Ctx) string {
		// limit by jwt token uid
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		uid := claims["uid"].(string)

		return uid
	},
	LimitReached: func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTooManyRequests, "too many requests")
	},
})
