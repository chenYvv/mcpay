package config

import (
	"bytes"
	"github.com/spf13/viper"
	"os"
)

func InitConfig(name string) {
	viper.SetConfigType("yml")
	config, err := os.ReadFile(name)
	if err != nil {
		panic("read config error:\n" + err.Error())
	}
	if err := viper.ReadConfig(bytes.NewBuffer(config)); err != nil {
		panic("read config error:\n" + err.Error())
	}
}
