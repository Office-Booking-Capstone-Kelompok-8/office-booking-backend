package config

import "time"

const (
	ACCESS_TOKEN_DURATION  = 15 * time.Minute
	REFRESH_TOKEN_DURATION = 14 * 24 * time.Hour
)
