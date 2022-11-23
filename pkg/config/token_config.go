package config

import "time"

const (
	ACCESS_TOKEN_DURATION  = 15 * time.Minute
	REFRESH_TOKEN_DURATION = 30 * 24 * time.Hour
)
