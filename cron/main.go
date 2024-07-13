package main

import "booking-app/cron/api"

func main() {
	cron := api.NewDailyTasksManager()
	if err := cron.Init(); err != nil {
		panic(err)
	}

	// 保持程序运行
	select {}
}
