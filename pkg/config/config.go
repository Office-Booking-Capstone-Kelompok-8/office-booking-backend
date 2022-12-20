package config

import "github.com/spf13/viper"

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
