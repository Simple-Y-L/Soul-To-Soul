package inits

import (
	"kratosItems/cmd/basic/config"

	"github.com/spf13/viper"
)

func Config() {
	viper.SetConfigFile("../../configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置信息失败" + err.Error())
	}
	if err := viper.Unmarshal(&config.Configs); err != nil {
		panic("序列化配置信息失败" + err.Error())
	}
}
