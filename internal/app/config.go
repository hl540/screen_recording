package app

import (
	"log"

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
		log.Printf("配置解析失败,err:%s", err.Error())
	}
}