package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func sendNotification() {
	fmt.Println("发送通知:", time.Now())
	// 这里实现你的通知发送逻辑，例如调用 API 或发送邮件
}

func SkipTestCron(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	// 设置调度规则为每个工作日的上午10点
	_, err := c.AddFunc("0 18 16 * * *", sendNotification)
	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return
	}

	c.Start()
	fmt.Println("调度程序已启动。每个工作日上午10点发送通知。")

	// 保持程序运行
	// select {}
}
