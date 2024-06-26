package core

import (
	"booking-app/micro-service/cluster/common/utils"

	"github.com/spf13/viper"
)

// 初始化配置
func InitConfig(filePath string) error {
	v := viper.New()
	v.SetConfigFile(filePath)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&Config); err != nil {
		return err
	}

	if Config.System.Server == "" {
		Config.System.Server = utils.GetLocalIP()
	}

	return nil
}
