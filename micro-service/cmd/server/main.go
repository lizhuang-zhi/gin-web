package main

import (
	"booking-app/micro-service/cluster/activity"
	"booking-app/micro-service/cluster/lobby"
	"booking-app/micro-service/core/app"
	"sync"
)

func main() {
	cluster := app.New("micro-service", "micro-service app")

	cluster.BeforeRun = func(app *app.App) error {
		// 启动服务
		var wg sync.WaitGroup

		activity.Start(&wg) // 启动活动服务
		lobby.Start(&wg)    // 启动大厅服务
		// 启动其他服务.....(后面整合)

		wg.Wait()

		return nil
	}

	cluster.Run()
}
