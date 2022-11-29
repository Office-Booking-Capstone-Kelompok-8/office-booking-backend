package config

import (
	"os"
	"time"
)

// Fiber config
//
//goland:noinspection ALL
const (
	APP_NAME             = "office-zone-api v0.1"
	SERVER_HEADER        = "office-zone-api"
	READ_TIMEOUT_SECONDS = 10
	SHUTDOWN_TIMEOUT     = 15
)

// JWT Access token and refresh token config
const (
	ACCESS_TOKEN_DURATION  = 15 * time.Minute
	REFRESH_TOKEN_DURATION = 14 * 24 * time.Hour
)

// OTP config
const (
	OTP_EXPIRATION_TIME = 15 * time.Minute
	OTP_LENGTH          = 6
	OTP_RESEND_TIME     = 1 * time.Minute
)

const (
	USER_ROLE  = 1
	ADMIN_ROLE = 2
)

func LoadConfig() map[string]string {
	env := make(map[string]string)

	env["DB_HOST"] = os.Getenv("DB_HOST")
	env["DB_PORT"] = os.Getenv("DB_PORT")
	env["DB_USER"] = os.Getenv("DB_USER")
	env["DB_PASS"] = os.Getenv("DB_PASS")
	env["DB_NAME"] = os.Getenv("DB_NAME")
	env["REDIS_HOST"] = os.Getenv("REDIS_HOST")
	env["REDIS_PORT"] = os.Getenv("REDIS_PORT")
	env["REDIS_PASS"] = os.Getenv("REDIS_PASS")
	env["REDIS_DB"] = os.Getenv("REDIS_DB")
	env["REFRESH_SECRET"] = os.Getenv("REFRESH_SECRET")
	env["ACCESS_SECRET"] = os.Getenv("ACCESS_SECRET")
	env["MAIL_DOMAIN"] = os.Getenv("MAIL_DOMAIN")
	env["MAIL_API_KEY"] = os.Getenv("MAIL_API_KEY")
	env["MAIL_SENDER"] = os.Getenv("MAIL_SENDER")
	env["MAIL_SENDER_NAME"] = os.Getenv("MAIL_SENDER_NAME")
	env["IMGKIT_PUBLIC_KEY"] = os.Getenv("IMGKIT_PUBLIC_KEY")
	env["IMGKIT_PRIVATE_KEY"] = os.Getenv("IMGKIT_PRIVATE_KEY")
	env["IMGKIT_URL_ENDPOINT"] = os.Getenv("IMGKIT_URL_ENDPOINT")
	env["PORT"] = os.Getenv("PORT")
	env["PREFORK"] = os.Getenv("PREFORK")

	return env
}
