package entity

type CachedToken struct {
	AccessID  string `json:"accessID"`
	RefreshID string `json:"refreshID"`
}

type CachedOTP struct {
	OTP    string `json:"otp"`
	Key    string `json:"tid"`
	UserID string `json:"uid"`
}
