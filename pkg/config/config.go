package config

import "github.com/spf13/viper"

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

func LoadConfig() (*viper.Viper, error) {
	conf := viper.New()
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath("./config")

	err := conf.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
