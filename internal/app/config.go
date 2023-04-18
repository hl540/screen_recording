package app

import (
	"github.com/spf13/viper"
)

var GlobalConfig *viper.Viper
var defaultConfigPath = "./"

func init() {
	GlobalConfig = viper.New()
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.SetConfigFile(defaultConfigPath + "conf.yaml")
	err := GlobalConfig.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

type Config struct{}
