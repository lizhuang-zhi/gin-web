package main

import (
	"booking-app/micro-service/cluster/activity"
	"booking-app/micro-service/cluster/common/config"
	"booking-app/micro-service/cluster/common/logger"
	"booking-app/micro-service/cluster/common/mongodb"
	"fmt"
)

func main() {
	logger.Info("start micro-service...")

	// 启动服务
	err := activity.Start()
	if err != nil {
		panic(fmt.Errorf("start server err:%v", err))
	}

	// 启动其他服务.....(后面整合)
}

func init() {
	// 读取配置
	config.LoadConfig()

	// 初始化日志
	logger.NewLogger()
	defer logger.GetLogger().Sync()

	// 初始化mongodb
	mongodb.NewMongoClient()
}
