package config

import (
	"os"
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
	env["PORT"] = os.Getenv("PORT")
	env["PREFORK"] = os.Getenv("PREFORK")

	return env
}
