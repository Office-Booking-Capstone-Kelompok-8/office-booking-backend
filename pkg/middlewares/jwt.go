package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"office-booking-backend/pkg/config"
	err2 "office-booking-backend/pkg/errors"
)

func NewJWTMiddleware(tokenSecret string, validator fiber.Handler) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(tokenSecret),
		ContextKey:     "user",
		SuccessHandler: validator,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err2.ErrInvalidToken.Error())
		},
	})
}

type AccessTokenValidator interface {
	CheckToken(ctx context.Context, token *jwt.MapClaims) (bool, error)
}

func ValidateAccessToken(validator AccessTokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		valid, err := validator.CheckToken(c.Context(), &claims)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}

		if !valid {
			return fiber.NewError(fiber.StatusUnauthorized, err2.ErrInvalidToken.Error())
		}

		return c.Next()
	}
}

func ValidateAdminAccessToken(validator AccessTokenValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if claims["role"] == float64(config.USER_ROLE) {
			return fiber.NewError(fiber.StatusForbidden, err2.ErrNoPermission.Error())
		}

		valid, err := validator.CheckToken(c.Context(), &claims)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}

		if !valid {
			return fiber.NewError(fiber.StatusUnauthorized, err2.ErrInvalidToken.Error())
		}

		return c.Next()
	}
}
