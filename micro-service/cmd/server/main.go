package main

import (
	"booking-app/micro-service/cluster/activity"
	"booking-app/micro-service/core/app"
	"fmt"
)

func main() {
	cluster := app.New("micro-service", "micro-service app")

	cluster.BeforeRun = func(app *app.App) error {
		// 启动服务
		err := activity.Start()
		if err != nil {
			panic(fmt.Errorf("start server err:%v", err))
		}

		// 启动其他服务.....(后面整合)

		return nil
	}

	cluster.Run()
}
