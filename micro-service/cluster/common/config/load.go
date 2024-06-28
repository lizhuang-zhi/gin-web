package config

import (
	"booking-app/micro-service/cluster/common"
	"booking-app/micro-service/cluster/common/utils"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 读取配置
func LoadConfig() {
	// 读取命令行参数
	configFilePath := pflag.StringP("config", "c", "./configs/local/config.yaml", "config file path")
	pflag.Parse()

	// 初始化配置
	err := InitConfig(*configFilePath)
	if err != nil {
		panic(fmt.Errorf("init config err:%v", err))
	}
}

// 初始化配置
func InitConfig(filePath string) error {
	v := viper.New()
	v.SetConfigFile(filePath)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&common.Config); err != nil {
		return err
	}

	if common.Config.System.Server == "" {
		common.Config.System.Server = utils.GetLocalIP()
	}

	return nil
}
