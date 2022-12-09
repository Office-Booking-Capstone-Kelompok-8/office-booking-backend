package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v4"
)

type Limiter struct {
	expiration time.Duration
}

func NewLimiter(expiration time.Duration) *Limiter {
	return &Limiter{
		expiration: expiration,
	}
}

func (l Limiter) ResetPasswordOTPLimitter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        1,
		Expiration: l.expiration,
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
}

func (l Limiter) VerifyEmailOTPLimitter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        1,
		Expiration: l.expiration,
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
}
