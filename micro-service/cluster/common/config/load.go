package config

import (
	"booking-app/micro-service/cluster/common"
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

	if err := v.Unmarshal(&common.Config); err != nil {
		return err
	}

	if common.Config.System.Server == "" {
		common.Config.System.Server = utils.GetLocalIP()
	}

	return nil
}
