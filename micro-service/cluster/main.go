package main

import (
	"booking-app/micro-service/cluster/activity"
	"booking-app/micro-service/cluster/common/core"
	"fmt"

	"github.com/spf13/pflag"
)

func main() {
	// 读取配置
	LoadConfig()

	// 启动服务
	activity := activity.NewServerInstance()
	err := activity.Start()
	if err != nil {
		panic(fmt.Errorf("start server err:%v", err))
	}

	// 启动其他服务.....(后面整合)
}

func LoadConfig() {
	// 读取命令行参数
	configFilePath := pflag.StringP("config", "c", "./configs/local/config.yaml", "config file path")
	pflag.Parse()

	// 初始化配置
	err := core.InitConfig(*configFilePath)
	if err != nil {
		panic(fmt.Errorf("init config err:%v", err))
	}
}
