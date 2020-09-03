package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(filename string, v interface{}) {
	config := viper.New()
	initConfigFromFiles(config, filename)

	config.Unmarshal(&v)
}

func initConfigFromFiles(config *viper.Viper, fileName string) {
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath("config")
	config.AddConfigPath("../config")

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file [%s]: %s \n", fileName, err))
	}
}
