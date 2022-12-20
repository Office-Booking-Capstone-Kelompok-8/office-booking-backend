package constant

import (
	"time"
)

// OTP config
const (
	RESET_PASSWORD_SUBJECT = "pass"
	VERIFY_EMAIL_SUBJECT   = "verify"
)

// Global config
const (
	USER_ROLE           = 1
	ADMIN_ROLE          = 2
	DEFAULT_USER_AVATAR = "https://ik.imagekit.io/fortyfour/default-image.jpg"
)

const (
	DATE_RESPONSE_FORMAT = "2006-01-02 15:04:05"
)

const (
	PAYMENT_EXPIRATION_TIME = 2 * 24 * time.Hour
)

const (
	PENDING_STATUS          = 1
	REJECTED_STATUS         = 2
	CANCELED_STATUS         = 3
	AWAITING_PAYMENT_STATUS = 4
	ACTIVE_STATUS           = 5
	COMPLETED_STATUS        = 6
)
