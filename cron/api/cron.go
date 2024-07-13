package api

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var DefaultDailyTasksManager *DailyTasksManager

// [每日定时]个人任务提醒
type DailyTasksManager struct {
	cron *cron.Cron
	stop chan struct{}
}

func NewDailyTasksManager() *DailyTasksManager {
	return &DailyTasksManager{
		cron: cron.New(),
		stop: make(chan struct{}),
	}
}

func (a *DailyTasksManager) Init() error {
	// 工作日早上10点提醒
	_, err := a.cron.AddFunc("04 18 * * *", func() {
		a.PushTask() // 推送个人任务提醒
	})
	if err != nil {
		return err
	}

	go func() {
		a.cron.Start()
		<-a.stop
		a.cron.Stop()
	}()

	return nil
}

// 推送个人任务提醒
func (a *DailyTasksManager) PushTask() {
	fmt.Println("PushTask start")

	// TODO: 推送核心逻辑
}

// 停止定时器
func (a *DailyTasksManager) Close() {
	close(a.stop)
}
