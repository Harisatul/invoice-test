package config

import (
	"github.com/spf13/viper"
	"os"
)

func newViper(name string) *viper.Viper {
	config := viper.New()

	config.SetConfigName(name)
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = os.Setenv("TZ", config.GetString("server.timezone"))
	if err != nil {
		panic(err)
	}

	return config
}
