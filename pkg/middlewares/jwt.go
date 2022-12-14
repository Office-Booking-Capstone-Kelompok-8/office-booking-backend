package middlewares

import (
	"context"
	"net/http"
	"office-booking-backend/pkg/constant"
	err2 "office-booking-backend/pkg/errors"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func NewJWTMiddleware(tokenSecret string, validator fiber.Handler) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(tokenSecret),
		ContextKey:     "user",
		SuccessHandler: validator,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return fiber.NewError(fiber.StatusUnauthorized, http.StatusText(fiber.StatusUnauthorized))
			}

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

		if claims["role"] == float64(constant.USER_ROLE) {
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

func EnforceValidEmail() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if claims["isVerified"] != true {
			return fiber.NewError(fiber.StatusForbidden, err2.ErrEmailNotVerified.Error())
		}

		return c.Next()
	}
}
