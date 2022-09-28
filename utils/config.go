package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver     string        `mapstructure:"DB_DRIVER"`
	DBSource     string        `mapstructure:"DB_SOURCE"`
	ServerAdress string        `mapstructure:"SERVER_ADDRESS"`
	SecretKey    string        `mapstructure:"SECRET_KEY"`
	KeyDuration  time.Duration `mapstructure:"KEY_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //json yml

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
