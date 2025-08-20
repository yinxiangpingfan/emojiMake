package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AI struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"ai"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(err)
	}
}
